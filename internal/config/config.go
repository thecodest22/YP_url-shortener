package config

// Решил поюзать pflag, чтобы сделать флаги в POSIX-стиле.
// Не понравилось, что стандартный "flag" не видит разницы
// между полным и сокращенным флагами и не дает собрать флаги в "-xyz".
import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// Константы нужны, т.к. ими пользуемся и при получении структуры конфигов с дефолтами,
// и в тестах, чтобы обеспечить соответствие.
const (
	shortenerCmdName         = "shortener"
	defaultLogLevel          = slog.LevelInfo
	defaultServerAddr        = "localhost:8080"
	defaultBaseURL           = "http://localhost:8080"
	defaultReadHeaderTimeout = 5 * time.Second
	defaultShutdownTimeout   = 5 * time.Second
)

type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type Config struct {
	Server   *ServerConfig
	LogLevel slog.Level
	BaseURL  string
}

var (
	cfg       Config
	fs        *pflag.FlagSet
	logLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

var (
	ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")
	ErrMissingURLHostname   = errors.New("missing URL hostname")
	ErrUnsupportedLogLevel  = errors.New("unsupported log level")
)

// defaultConfig отдает экземпляр структуры корневого конфига с дефолтами.
// Потом поля этой структуры можно будет переопределять после чтения конфигов из env и/или из CLI.
func defaultConfig() Config {
	return Config{
		Server: &ServerConfig{
			Addr:              defaultServerAddr,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ShutdownTimeout:   defaultShutdownTimeout,
		},
		LogLevel: defaultLogLevel,
		BaseURL:  defaultBaseURL,
	}
}

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
		&cfg.Server.Addr,
		"server-address",
		"a",
		defaultServerAddr,
		"HTTP server listen address",
	)

	// Регистрируем аргументы с кастомными обработчиками.
	// Они не ставит дефолты, но мы заранее получаем дефолтную конфигурацию, прежде чем она сюда попадет.
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

	flagSet.FuncP(
		"log-level",
		"l",
		"log level",
		func(value string) error {
			lvl, exists := logLevels[strings.ToLower(value)]
			if !exists {
				return fmt.Errorf(
					"%w: unsupported level %q, allowed levels: 'debug', 'info', 'warn', 'error'",
					ErrUnsupportedLogLevel,
					value,
				)
			}

			cfg.LogLevel = lvl
			return nil
		},
	)

	return flagSet
}

// Заюзал init просто ради экспириенса
func init() {
	cfg = defaultConfig()
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
