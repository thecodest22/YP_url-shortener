package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
)

// FIXME убрать, когда будем тянуть конфиги извне
const (
	logLevel = slog.LevelDebug
	addr     = "localhost:8080"
)

// Run - оркестрирует запуск приложения.
// Выполняет подготовительные действия (тянет конфиги, подключает репы и т.д.),
// инициализирует и запускает сервак.
func Run() error {
	// TODO объявление контекста

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	// TODO сбор конфигов

	// TODO подключение репозиториев

	srv := httpserver.New(addr, handler.NewRouter(logger))
	logger.Info(
		"starting http server",
		"addr", addr,
		"log_level", logLevel,
	)

	if err := srv.Run(); err != nil {
		// Тут не будем логировать. Только врапаем ошибку,
		// а все вместе залогируется в main.
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}
