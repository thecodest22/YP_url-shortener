package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var (
	// logLevels - мапа для сопоставления значений флага CLI с константами пакета логирования.
	logLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

var (
	ErrUnsupportedLogLevel = errors.New("unsupported log level")
)

// Config конфигурация для логгера.
type Config struct {
	Level string
}

func New(cfg Config) (*slog.Logger, error) {
	lvl, ok := logLevels[strings.ToLower(cfg.Level)]
	if !ok {
		return nil, fmt.Errorf(
			"log level %q: %w",
			cfg.Level,
			ErrUnsupportedLogLevel,
		)
	}

	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: lvl},
		),
	), nil
}
