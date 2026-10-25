package httpserver

import (
	"net/http"
	"time"
)

// FIXME убрать это, когда начнем тянуть конфиги извне
const readHeaderTimeout = 5 * time.Second

// Server - http-сервер приложения.
type Server struct {
	httpServer *http.Server
}

// New - инициализирует http-сервер с ручками.
// TODO принимать не только адрес, а сразу структуру конфигов (адрес, таймауты...)
func New(addr string, h http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           h,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

// Run - запускает http-сервер.
func (s *Server) Run() error {
	// TODO принимать тут контекст и запустить тут ListenAndServe как горутину
	//  с каналами, ошибками и т.д.
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
