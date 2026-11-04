package config

// Решил поюзать pflag, чтобы сделать флаги в POSIX-стиле.
// Не понравилось, что стандартный "flag" не видит разницы
// между полным и сокращенным флагами и не дает собрать флаги в "-xyz".
import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/pflag"
)

const (
	shortenerCmdName = "shortener"
	defaultBaseURL   = "http://localhost:8080"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

var (
	cfg Config
	fs  *pflag.FlagSet
)

var (
	ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")
	ErrMissingURLHostname   = errors.New("missing URL hostname")
)

// validateBaseURL валидирует "--base-url" из командной строки.
func validateBaseURL(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}

	// Тут врапнем, чтобы в сообщении содержалось, какой протокол был на входе
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: %q", ErrUnsupportedURLScheme, u.Scheme)
	}

	if u.Hostname() == "" {
		return ErrMissingURLHostname
	}

	return nil
}

// newFlagSet объявляет флаги для парсинга.
// Пришлось вынести в отдельную функцию из init, чтобы было удобно тестировать.
// Иначе после парсинга аргументов изменяется общее состояние, дефолты слетают.
func newFlagSet(cfg *Config) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(shortenerCmdName, pflag.ContinueOnError)

	// Это будем расценивать, как обычный строковый флаг, без кастомных валидаций.
	// Тут дефолтное значение при отсутствии аргумента будет задаваться автоматически.
	flagSet.StringVarP(
		&cfg.ServerAddr,
		"server-address",
		"a",
		"localhost:8080",
		"HTTP server listen address",
	)

	// Руками ставим дефолты, т.к. при применении кастомной функции-валидатора
	// дефолты автоматически ставиться не будут!
	cfg.BaseURL = defaultBaseURL

	// Регистрируем аргумент с кастомным обработчиком. Он не ставит дефолты.
	// До функции-валидатора дело доходит, только если аргумент передан и распарсен.
	// Здесь мы должны сами проверить формат, значение и т.д.
	flagSet.FuncP(
		"base-url",
		"b",
		"base URL for shortened links",
		func(value string) error {
			// Если здесь получаем ошибку, и возвращаем наверх, pflag оборачивает ее
			// в *pflag.InvalidValueError, добавляя сведения о флаге и значении.
			// Затем fs.Parse(args) возвращает эту обернутую ошибку,
			// а наш Parse возвращает Config{}, err. ИТОГО: удобные ошибки и удобно тестировать.
			if err := validateBaseURL(value); err != nil {
				return err
			}

			// Сохраняем только успешно проверенное значение.
			cfg.BaseURL = value
			return nil
		},
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
