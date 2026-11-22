package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"

	"github.com/thecodest22/YP_url-shortener/internal/config"
	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/logger"
	"github.com/thecodest22/YP_url-shortener/internal/service"
	"github.com/thecodest22/YP_url-shortener/internal/storage/memory"
)

// Run - оркестрирует запуск приложения.
// Выполняет подготовительные действия (тянет конфиги, подключает репы и т.д.),
// инициализирует и запускает сервак.
func Run() error {
	cfg, err := config.Load(os.Args[1:])

	// Если это был флаг --help, это тоже считается "ошибкой", но на самом деле
	// pflag вывел справку. Тогда просто завершаем приложение.
	// Если же реально ошибка парсинга, возвращаем ошибку.
	if err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil
		}
		return err
	}

	// Объявляем контекст, слушающий сигналы остановки.
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)
	defer stop()
	// Регистрация функции для немедленного вызова в момент отмены контекста.
	// Нужно, чтобы при отмене контекста сразу вызвался stop, который отменит
	// перехват сигналов. Иначе повторное нажатие Ctrl+C так и не будет
	// принудительно завершать процесс, пока Run не доберется до defer stop().
	context.AfterFunc(ctx, stop)

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	st := memory.New()
	svc := service.NewURLShortener(st)
	router, err := handler.NewRouter(log, svc, cfg.Handler)
	if err != nil {
		return fmt.Errorf("init router: %w", err)
	}
	srv := httpserver.New(cfg.Server, router)

	log.Info(
		"starting http server",
		"addr", cfg.Server.Addr,
		"log_level", cfg.Logger.Level,
	)

	if err := srv.Run(ctx); err != nil {
		// Тут не будем логировать. Только врапаем ошибку,
		// а все вместе залогируется в main.
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}
