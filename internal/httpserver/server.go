package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/thecodest22/YP_url-shortener/internal/config"
)

// Server - http-сервер приложения.
type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

// New - инициализирует http-сервер с ручками.
// Принимаем на вход хэндлер (будем передавать "корневой" роутер), и тогда сервер будет
// слушать запросы и передавать в функции этого хэндлера (их регистрируем соппсна на корневом роутере).
func New(cfg config.ServerConfig, h http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.Addr,
			Handler:           h,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

// Run - запускает http-сервер.
// Принимает контекст, слушающий сигналы остановки.
// Возвращает ошибку работы сервера (или nil - если это ожидаемая ошибка остановки
// сервера, т.е., по сути, не ошибочное состояние).
func (s *Server) Run(ctx context.Context) error {
	// Чтобы контекст Request'а слушал сигналы отмены и передавал в вызываемые
	// функции (например, функции репозитория), нужно переопределить BaseContext
	// на httpServer, иначе контексты Request'ов будут "наследоваться" от обычного
	// context.Background(), и связи с нашим NotifyContext не будет.
	s.httpServer.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	// Объявляем "errgroup" - это координатор горутин. Гарантирует корректное
	// завершение зарегистрированных в нем горутин.
	// Если любая из горутин вернет error, то остальные получат сигнал остановки
	// через общий контекст. Благодаря ему не нужно делать select'ы и вручную
	// обрабатывать ошибки, возвращаемые горутинами через канал.
	g, gCtx := errgroup.WithContext(ctx) // "g" - указатель на объект Group.

	// ------ Регистрируем горутины ------
	// 1) Горутина запуска сервера. Нужна, чтобы ListenAndServe() выполнялся
	// асинхронно и не блокировал рантайм в main().
	g.Go(func() error {
		err := s.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	// 2) Горутина остановки сервака при отмене или ошибке другой горутины.
	g.Go(func() error {
		<-gCtx.Done() // В случае ошибки другой горутины здесь будет как будто отмена.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	})
	// ------------

	// Блокирует и дожидается завершения всех зареганных горутин.
	// Возвращает ошибку из упавшей горутины либо nil.
	return g.Wait()
}
