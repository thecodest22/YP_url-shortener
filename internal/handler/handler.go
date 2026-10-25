package handler

import (
	"log/slog"
	"net/http"
)

// handler - структура для корневого хэндлера. Здесь общие поля, нужные для
// методов-хэндлеров, объявленных на структуре.
type handler struct {
	logger *slog.Logger
}

// NewRouter - создает структуру.
func NewRouter(logger *slog.Logger) http.Handler {
	h := &handler{logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}
