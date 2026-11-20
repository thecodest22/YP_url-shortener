package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

// TestApplyCLI тестирует функцию применения флагов командной строки к переданному конфигу.
func TestApplyCLI(t *testing.T) {
	// На входе будет конфиг с уже какими-то значениями (не важно, откуда они пришли).
	initial := rawConfig{
		Server: rawServerConfig{
			Addr:              "localhost:9090",
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ShutdownTimeout:   defaultShutdownTimeout,
		},
		Logger: rawLoggerConfig{
			Level: "debug",
		},
		Handler: rawHandlerConfig{
			BaseURL: "https://ya.ru",
		},
	}

	tests := []struct {
		name string
		args []string
		want rawConfig
	}{
		{
			name: "positive: absent flags preserve existing values",
			args: []string{},
			want: initial,
		},
		{
			name: "positive: passed flags override existing values",
			args: []string{"-l", "warn"},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:9090",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "warn",
				},
				Handler: rawHandlerConfig{
					BaseURL: "https://ya.ru",
				},
			},
		},
		{
			name: "positive: passed flags override existing values",
			args: []string{"-a", "localhost:8888"},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:8888",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "debug",
				},
				Handler: rawHandlerConfig{
					BaseURL: "https://ya.ru",
				},
			},
		},
		{
			name: "positive: passed flags override existing values",
			args: []string{"-b", "https://google.com"},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:9090",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "debug",
				},
				Handler: rawHandlerConfig{
					BaseURL: "https://google.com",
				},
			},
		},
		{
			name: "positive: values override 'as is' without validation",
			args: []string{"--log-level", "UNSUPPORTED"},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:9090",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "UNSUPPORTED",
				},
				Handler: rawHandlerConfig{
					BaseURL: "https://ya.ru",
				},
			},
		},
		{
			name: "positive: values override 'as is' without validation",
			args: []string{"--base-url", ""},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:9090",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "debug",
				},
				Handler: rawHandlerConfig{
					BaseURL: "",
				},
			},
		},
		{
			name: "positive: values override only if passed explicitly",
			// Суть в том, что у нас INFO зареган во флагсете как дефолтный.
			// Но примениться он должен только тогда, когда передан явно.
			args: []string{"-l", defaultLogLevel.String()},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:9090",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "INFO",
				},
				Handler: rawHandlerConfig{
					BaseURL: "https://ya.ru",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			got := initial

			var cli cliOptions
			fs := newFlagSet(&cli)
			// Если почему-то отвалился парсинг аргументов, дальше не продолжаем тест.
			require.NoError(t, fs.Parse(test.args))

			// Применяем флаги к конфигу на входе.
			applyCLI(&got, cli, fs)

			// ------ Проверки ------
			assert.Equal(t, test.want, got)
		})
	}
}
