package config

import (
	"log/slog"
	"time"
)

// Константы нужны, т.к. ими пользуемся и при получении структуры конфигов с дефолтами,
// и в тестах, чтобы обеспечить соответствие.
const (
	defaultLogLevel          = slog.LevelInfo
	defaultServerAddr        = "localhost:8080"
	defaultBaseURL           = "http://localhost:8080"
	defaultReadHeaderTimeout = 5 * time.Second
	defaultShutdownTimeout   = 5 * time.Second
)

// defaultRawConfig создает начальную конфигурацию в виде rawConfig
// с дефолтами уровня кода.
func defaultRawConfig() rawConfig {
	return rawConfig{
		Server: rawServerConfig{
			Addr:              defaultServerAddr,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ShutdownTimeout:   defaultShutdownTimeout,
		},
		Logger: rawLoggerConfig{
			Level: defaultLogLevel.String(),
		},
		Handler: rawHandlerConfig{
			BaseURL: defaultBaseURL,
		},
	}
}
