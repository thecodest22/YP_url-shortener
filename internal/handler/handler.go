package handler

import (
	"context"
	"log/slog"
	"net/http"
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
	Resolve(ctx context.Context, id string) (string, error)
}

// NewRouter - создает новый корневой роутер с нужными хэндлерами.
func NewRouter(logger *slog.Logger, svc urlShortener, baseURL string) http.Handler {
	h := &handler{
		logger:  logger,
		svc:     svc,
		baseURL: baseURL,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("POST /{$}", h.shorten) // "/{$}" означает не слушать тут запросы по каким-то другим маршрутам после слэша
	return mux
}
