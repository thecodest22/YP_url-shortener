package config

// Решил поюзать pflag, чтобы сделать флаги в POSIX-стиле.
// Не понравилось, что стандартный "flag" не видит разницы
// между полным и сокращенным флагами и не дает собрать флаги в "-xyz".
import (
	"github.com/spf13/pflag"
)

const shortenerCmdName = "shortener"

type Config struct {
	ServerAddr string
	BaseURL    string
}

var (
	cfg Config
	fs  *pflag.FlagSet
)

// newFlagSet объявляет флаги для парсинга.
// Пришлось вынести в отдельную функцию, чтобы было удобно тестировать.
// Иначе после парсинга аргументов изменяется общее состояние, дефолты слетают.
func newFlagSet(cfg *Config) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(shortenerCmdName, pflag.ContinueOnError)

	flagSet.StringVarP(
		&cfg.ServerAddr,
		"server-address",
		"a",
		"localhost:8080",
		"HTTP server listen address",
	)
	flagSet.StringVarP(
		&cfg.BaseURL,
		"base-url",
		"b",
		"http://localhost:8080",
		"base URL for shortened links",
	)

	return flagSet
}

// Заюзал init просто ради экспириенса
func init() {
	// Присваиваем результат в переменную пакета. Затем в тестах будем делать
	// то же самое, чтобы реинициализировать флагсет с его дефолтами.
	fs = newFlagSet(&cfg)
}

// Parse разбирает аргументы командной строки и кладет в структуру конфига.
func Parse(args []string) (Config, error) {
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
