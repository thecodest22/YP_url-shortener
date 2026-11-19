package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Config конфиг хэндлера.
type Config struct {
	BaseURL string
}

// handler - структура для использования в корневом роутере.
// Здесь общие поля, нужные для методов-хэндлеров, объявленных на роутере.
type handler struct {
	logger  *slog.Logger
	svc     urlShortener
	baseURL string
}

type urlShortener interface {
	// Shorten - сократить URL и вернуть его ID для последующего резолва.
	Shorten(ctx context.Context, original string) (id string, err error)

	// Resolve - вернуть оригинальный URL по его ID.
	Resolve(ctx context.Context, id string) (url string, err error)
}

// NewRouter - создает новый корневой роутер с нужными хэндлерами.
func NewRouter(
	logger *slog.Logger,
	svc urlShortener,
	cfg Config,
) (http.Handler, error) {
	if err := validateBaseURL(cfg.BaseURL); err != nil {
		return nil, fmt.Errorf("base URL: %w", err)
	}

	h := &handler{
		logger:  logger,
		svc:     svc,
		baseURL: cfg.BaseURL,
	}

	mux := chi.NewRouter()
	mux.Use(middleware.GetHead) // Добавляет обработку HEAD-запросов
	mux.Get("/healthz", h.health)
	mux.Post("/", h.shorten)
	mux.Get("/{id}", h.resolve)
	return mux, nil
}
