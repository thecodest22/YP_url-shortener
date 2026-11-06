package config

import (
	"log/slog"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCfg Config
		wantErr testutil.ErrorCheck
	}{
		{
			name:    "positive: no args passed",
			args:    []string{},
			wantCfg: defaultConfig(),
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"-l", "debug"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelDebug,
				BaseURL:  defaultBaseURL,
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"--log-level", "DEBUG"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelDebug,
				BaseURL:  defaultBaseURL,
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"-l", "warn"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelWarn,
				BaseURL:  defaultBaseURL,
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"-l", "error"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelError,
				BaseURL:  defaultBaseURL,
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with full name",
			args: []string{"--base-url", "http://foo.com"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelInfo,
				BaseURL:  "http://foo.com",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with shorthand",
			args: []string{"-b", "http://foo.com"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelInfo,
				BaseURL:  "http://foo.com",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelInfo,
				BaseURL:  "http://ya.ru",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "--base-url", "http://ya.ru"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelInfo,
				BaseURL:  "http://ya.ru",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"--base-url", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				Server: &ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				LogLevel: slog.LevelInfo,
				BaseURL:  "http://ya.ru",
			},
			wantErr: nil,
		},
		{
			name:    "positive: ignorable position args passed",
			args:    []string{"foo", "bar"},
			wantCfg: defaultConfig(),
			wantErr: nil,
		},
		{
			name:    "negative: unknown args passed",
			args:    []string{"--foo", "bar"},
			wantCfg: Config{},
			wantErr: testutil.ErrorAs[*pflag.NotExistError], // Тут функция только инициализируется из дженерика!
		},
		{
			name:    "negative: invalid base-url: not a URL",
			args:    []string{"-b", "bar"},
			wantCfg: Config{},
			wantErr: testutil.ErrorAs[*pflag.InvalidValueError],
		},
		{
			name:    "negative: invalid log level",
			args:    []string{"--log-level", "critical"},
			wantCfg: Config{},
			wantErr: testutil.ErrorIs(ErrUnsupportedLogLevel), // Тут возвращается обычная функция
		},
		{
			name:    "negative: invalid base-url: unsupported scheme",
			args:    []string{"-b", "ftp://foo.com"},
			wantCfg: Config{},
			wantErr: testutil.ErrorIs(ErrUnsupportedURLScheme),
		},
		{
			name:    "negative: invalid base-url: missing host name",
			args:    []string{"-b", "http:///some-path"},
			wantCfg: Config{},
			wantErr: testutil.ErrorIs(ErrMissingURLHostname),
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			// Т.к. мы работаем с переменными пакета,
			// нужно сначала сохранить состояние, которое было до теста:
			//savedCfg, savedFS := cfg, fs

			// Регистрируем восстановление с помощью специальной функции.
			// Она выполнится после завершения теста (даже если упадет на require).
			//t.Cleanup(func() {
			//	cfg = savedCfg
			//	fs = savedFS
			//})

			// Приводим переменные пакета в дефолтное состояние для текущего теста.
			//cfg = defaultConfig()
			//fs = newFlagSet(&cfg)

			gotCfg, err := Parse(test.args)

			if test.wantErr == nil {
				// Если должно быть без ошибок, то проверяем стандартной функцией
				require.NoError(t, err, "expected no error")
			} else {
				// Если должна быть ошибка, вызываем проверку кастомной функцией, которую сами определили
				test.wantErr(t, err)
			}
			assert.Equal(t, test.wantCfg, gotCfg, "got incorrect configuration")
		})
	}
}
