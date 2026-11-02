package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
func NewRouter(logger *slog.Logger, svc urlShortener, baseURL string) http.Handler {
	h := &handler{
		logger:  logger,
		svc:     svc,
		baseURL: baseURL,
	}

	mux := chi.NewRouter()
	mux.Get("/healthz", h.health)
	mux.Post("/", h.shorten)
	mux.Get("/{id}", h.resolve)
	return mux
}
