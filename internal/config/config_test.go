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
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: slog.LevelDebug,
				},
				Handler: HandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"--log-level", "DEBUG"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: slog.LevelDebug,
				},
				Handler: HandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"-l", "warn"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: slog.LevelWarn,
				},
				Handler: HandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: valid log level",
			args: []string{"-l", "error"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: slog.LevelError,
				},
				Handler: HandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with full name",
			args: []string{"--base-url", "http://foo.com"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "http://foo.com",
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with shorthand",
			args: []string{"-b", "http://foo.com"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "http://foo.com",
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "http://ya.ru",
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "--base-url", "http://ya.ru"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "http://ya.ru",
				},
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"--base-url", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "http://ya.ru",
				},
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
			name: "positive: not a URL", // Парсер не отвечает за валидацию, ее выполняет потребитель
			args: []string{"-b", "bar"},
			wantCfg: Config{
				Server: ServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: LoggerConfig{
					Level: defaultLogLevel,
				},
				Handler: HandlerConfig{
					BaseURL: "bar",
				},
			},
			wantErr: nil,
		},
		{
			name:    "negative: unknown args passed",
			args:    []string{"--foo", "bar"},
			wantCfg: Config{},
			wantErr: testutil.ErrorAs[*pflag.NotExistError], // Тут функция только инициализируется из дженерика!
		},
		{
			name:    "negative: invalid log level",
			args:    []string{"--log-level", "critical"},
			wantCfg: Config{},
			wantErr: testutil.ErrorIs(ErrUnsupportedLogLevel), // Тут возвращается обычная функция
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
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
