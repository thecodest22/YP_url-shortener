package config

// Решил поюзать pflag, чтобы сделать флаги в POSIX-стиле.
// Не понравилось, что стандартный "flag" не видит разницы
// между полным и сокращенным флагами и не дает собрать флаги в "-xyz".
import (
	"errors"
	"fmt"
	"log/slog"
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

// ServerConfig - конфиги сервера.
type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type LoggerConfig struct {
	Level slog.Level
}

type HandlerConfig struct {
	BaseURL string
}

// Config - корневой конфиг приложения.
type Config struct {
	Server  ServerConfig
	Logger  LoggerConfig
	Handler HandlerConfig
}

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

// defaultConfig отдает экземпляр структуры корневого конфига с дефолтами.
// Потом поля этой структуры можно будет переопределять после чтения конфигов из env и/или из CLI.
func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr:              defaultServerAddr,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ShutdownTimeout:   defaultShutdownTimeout,
		},
		Logger: LoggerConfig{
			Level: defaultLogLevel,
		},
		Handler: HandlerConfig{
			BaseURL: defaultBaseURL,
		},
	}
}

// newFlagSet объявляет флаги для парсинга.
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

	flagSet.StringVarP(
		&cfg.Handler.BaseURL,
		"base-url",
		"b",
		defaultBaseURL,
		"base URL for shortened links",
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

			cfg.Logger.Level = lvl
			return nil
		},
	)

	return flagSet
}

// Parse разбирает аргументы командной строки и кладет в структуру конфига.
func Parse(args []string) (Config, error) {
	cfg := defaultConfig()
	fs := newFlagSet(&cfg)

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
