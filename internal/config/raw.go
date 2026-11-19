package config

import (
	"time"

	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/logger"
)

type rawServerConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type rawLoggerConfig struct {
	Level string
}

type rawHandlerConfig struct {
	BaseURL string
}

// rawConfig - структура конфигов до окончательного преобразования.
// Сначала инициализиурем ее дефолтами проги, а потом применим значения
// из разных стадий сбора конфигов - файл, флаги, окружение.
// Здесь можем объявлять типы как базовые (string, int и т.д.), чтобы сначала выбрать значение
// из источника согласно приоритетам, а в конце провалидировать и привести к типу.
type rawConfig struct {
	Server  rawServerConfig
	Logger  rawLoggerConfig
	Handler rawHandlerConfig
}

// toConfig возвращает Config со значениями rawConfig
// после всех этапов применения значений конфигов из разных источников.
func (raw rawConfig) toConfig() Config {
	return Config{
		Server: httpserver.Config{
			Addr:              raw.Server.Addr,
			ReadHeaderTimeout: raw.Server.ReadHeaderTimeout,
			ShutdownTimeout:   raw.Server.ShutdownTimeout,
		},
		Logger: logger.Config{
			Level: raw.Logger.Level,
		},
		Handler: handler.Config{
			BaseURL: raw.Handler.BaseURL,
		},
	}
}
