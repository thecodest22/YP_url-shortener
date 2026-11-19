package config

import (
	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/logger"
)

// Config - корневой конфиг приложения.
type Config struct {
	Server  httpserver.Config
	Logger  logger.Config
	Handler handler.Config
}

// Load собирает конфигурацию приложения.
// Управляет флоу сбора конфигов из разных источников и их применением в соответствии с приоритетами.
func Load(args []string) (Config, error) {
	// Получаем дефолтные настройки.
	raw := defaultRawConfig()

	// Затем сбор флагов из командной строки. Сами конфиги применим потом,
	// но сами флаги нужно получить заранее, т.к. там может быть расположение файла переменных.
	var cli cliOptions
	fs := newFlagSet(&cli)

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	// TODO Если флаг файла передан, то читаем файл и применяем значения

	// Применяем флаги из CLI.
	applyCLI(&raw, cli, fs)

	// TODO Применяем переменные окружения

	return raw.toConfig(), nil
}
