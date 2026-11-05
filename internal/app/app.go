package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/thecodest22/YP_url-shortener/internal/config"
	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/service"
	"github.com/thecodest22/YP_url-shortener/internal/storage/memory"
)

// Run - оркестрирует запуск приложения.
// Выполняет подготовительные действия (тянет конфиги, подключает репы и т.д.),
// инициализирует и запускает сервак.
func Run() error {
	cfg, err := config.Parse(os.Args[1:])

	// Если это был флаг --help, это тоже считается "ошибкой", но на самом деле
	// pflag вывел справку, хотя это считается ошибкой. Тогда просто завершаем приложение.
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	st := memory.New()
	svc := service.NewURLShortener(st)
	srv := httpserver.New(cfg.ServerAddr, handler.NewRouter(logger, svc, cfg.BaseURL))

	// TODO отдельно файлик с созданием логгера на основе конфигов
	logger.Info(
		"starting http server",
		"addr", cfg.ServerAddr,
		"log_level", cfg.LogLevel,
	)

	if err := srv.Run(ctx); err != nil {
		// Тут не будем логировать. Только врапаем ошибку,
		// а все вместе залогируется в main.
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}
