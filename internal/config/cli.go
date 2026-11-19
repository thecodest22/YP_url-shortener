package config

// Решил поюзать pflag, чтобы сделать флаги в POSIX-стиле.
// Не понравилось, что стандартный "flag" не видит разницы
// между полным и сокращенным флагами и не дает собрать флаги в "-xyz".
import "github.com/spf13/pflag"

// Имя группы команд
const (
	shortenerCmdName   = "shortener"
	configFlagName     = "config"
	serverAddrFlagName = "server-address"
	baseURLFlagName    = "base-url"
	logLevelFlagName   = "log-level"
)

// cliOptions хранит значения из флагов "как есть".
type cliOptions struct {
	ConfigPath string
	ServerAddr string
	BaseURL    string
	LogLevel   string
}

// newFlagSet создает набор флагов для группы команд и регистрирует флаги для парсинга.
func newFlagSet(cli *cliOptions) *pflag.FlagSet {
	// Объявляем флагсет для группы команд.
	// Поведение при обнаружении ошибки - pflag.ContinueOnError.
	// NOTE: Несмотря на название оно прекращает выполнение парсинга.
	fs := pflag.NewFlagSet(shortenerCmdName, pflag.ContinueOnError)

	// Регистрируем у флагсета флаги. StringVarP - обычный строковый флаг, без кастомных валидаций.
	// StringVarP при вызове сразу задает дефолтное значение для поля (по адресу).
	fs.StringVarP(
		&cli.ConfigPath,
		configFlagName,
		"c",
		"",
		"path to YAML configuration file",
	)

	fs.StringVarP(
		&cli.ServerAddr,
		serverAddrFlagName,
		"a",
		defaultServerAddr,
		"HTTP server listen address",
	)

	fs.StringVarP(
		&cli.BaseURL,
		baseURLFlagName,
		"b",
		defaultBaseURL,
		"base URL for shortened links",
	)

	fs.StringVarP(
		&cli.LogLevel,
		logLevelFlagName,
		"l",
		defaultLogLevel.String(),
		"log level",
	)

	return fs
}

// applyCLI применяет к конфигу значения, если они действительно были переданы через флаги.
func applyCLI(raw *rawConfig, cli cliOptions, fs *pflag.FlagSet) {
	if fs.Changed(serverAddrFlagName) {
		raw.Server.Addr = cli.ServerAddr
	}

	if fs.Changed(baseURLFlagName) {
		raw.Handler.BaseURL = cli.BaseURL
	}

	if fs.Changed(logLevelFlagName) {
		raw.Logger.Level = cli.LogLevel
	}
}
