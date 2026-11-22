package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setEnvs - хелпер, накидывает переменные среды на время теста.
func setEnvs(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		// Setenv запоминает исходное состояние и автоматически
		// восстановит его после теста.
		t.Setenv(k, v)
	}
}

// stashEnvs - хелпер, перед тестом очищает среду от переменных,
// чтобы они не конфликтовали с тестируемыми значениями.
// В конце выполнения теста восстанавливает исходные значения.
func stashEnvs(t *testing.T, keys []string) {
	t.Helper()

	// Setenv запоминает исходное состояние, но синтаксически функции нужно
	// новое значение, поэтому кидаем просто пустую строку.
	// Затем Unsetenv вообще удаляет переменную, а после теста Setenv вернет ее на место.
	for _, key := range keys {
		t.Setenv(key, "")
		require.NoError(t, os.Unsetenv(key))
	}
}

func TestReadEnv(t *testing.T) {
	// На входе будет конфиг с уже какими-то значениями (не важно, откуда они пришли).
	initial := defaultRawConfig()

	tests := []struct {
		name string
		envs map[string]string
		want rawConfig
	}{
		{
			name: "positive: missing values preserve existing values",
			envs: map[string]string{},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: defaultLogLevel.String(),
				},
				Handler: rawHandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
		},
		{
			name: "positive: missing values preserve existing values",
			envs: map[string]string{
				"LOG_LEVEL": "warn",
			},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "warn",
				},
				Handler: rawHandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
		},
		{
			name: "positive: empty strings passed 'as is'",
			envs: map[string]string{
				"LOG_LEVEL": "",
			},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "",
				},
				Handler: rawHandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
		},
		{
			name: "positive: full set of values",
			envs: map[string]string{
				"LOG_LEVEL":      "error",
				"SERVER_ADDRESS": "localhost:8888",
				"BASE_URL":       "http://test.com:9090",
			},
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:8888",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "error",
				},
				Handler: rawHandlerConfig{
					BaseURL: "http://test.com:9090",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d %s", i+1, test.name), func(t *testing.T) {
			got := initial

			// Прячем переменные среды и накидываем нужные.
			stashEnvs(
				t,
				[]string{
					"LOG_LEVEL",
					"SERVER_ADDRESS",
					"BASE_URL",
				},
			)
			setEnvs(t, test.envs)

			err := readEnv(&got)

			// ------ Проверки ------
			require.NoError(t, err, "failed to read env")
			assert.Equal(t, test.want, got, "unexpected config after reading env variables")
		})
	}
}
