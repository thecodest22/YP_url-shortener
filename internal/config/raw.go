package config

import (
	"time"

	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/logger"
)

type rawServerConfig struct {
	Addr              string        `yaml:"address"`
	ReadHeaderTimeout time.Duration `yaml:"-"` // Не читаем из ямлика
	ShutdownTimeout   time.Duration `yaml:"-"`
}

type rawLoggerConfig struct {
	Level string `yaml:"level"`
}

type rawHandlerConfig struct {
	BaseURL string `yaml:"base_url"`
}

// rawConfig - структура конфигов до окончательного преобразования.
// Сначала инициализиурем ее дефолтами проги, а потом применим значения
// из разных стадий сбора конфигов - файл, флаги, окружение.
// Здесь можем объявлять типы как базовые (string, int и т.д.), чтобы сначала выбрать значение
// из источника согласно приоритетам, а в конце провалидировать и привести к типу.
type rawConfig struct {
	Server  rawServerConfig  `yaml:"server"`
	Logger  rawLoggerConfig  `yaml:"logger"`
	Handler rawHandlerConfig `yaml:"handler"`
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
