package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	// На входе будет конфиг с уже какими-то значениями (не важно, откуда они пришли).
	initial := defaultRawConfig()

	tests := []struct {
		name        string
		yamlContent []byte
		want        rawConfig
	}{
		{
			name: "positive: missing values preserve existing values",
			yamlContent: []byte(`
                logger:
                  level: DEBUG
            `),
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              defaultServerAddr,
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "DEBUG",
				},
				Handler: rawHandlerConfig{
					BaseURL: defaultBaseURL,
				},
			},
		},
		{
			name: "positive: null values preserve existing values",
			yamlContent: []byte(`
                server:
                  address:
            `),
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
			name:        "positive: empty mapping preserves existing values",
			yamlContent: []byte("{}"), // Так создастся валидный ямлик, но он будет пустой
			want:        initial,
		},
		{
			name: "positive: empty empty string passed 'as is'",
			yamlContent: []byte(`
                server:
                  address: ""
            `),
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "",
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
			name: "positive: full set of values",
			yamlContent: []byte(`
                server:
                  address: "localhost:8888"
                logger:
                  level: warn
                handler:
                  base_url: http://test.com:9090
            `),
			want: rawConfig{
				Server: rawServerConfig{
					Addr:              "localhost:8888",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: rawLoggerConfig{
					Level: "warn",
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

			// Пишем файл во временную директорию TempDir ее сам создаст и удалит после теста.
			path := filepath.Join(t.TempDir(), "config.yaml")
			err := os.WriteFile(path, test.yamlContent, 0600)

			// Если отвалилась запись, дальше смысла нет тестить.
			require.NoError(t, err, "failed to write test config file %q", path)

			err = readFile(path, &got)

			// ------ Проверки ------
			require.NoError(t, err, "failed to read test config file %q", path)
			assert.Equal(t, test.want, got, "unexpected config after reading test config file")
		})
	}
}

func TestReadInvalidFile(t *testing.T) {
	initial := defaultRawConfig()

	path := filepath.Join(t.TempDir(), "config.yaml")
	err := os.WriteFile(
		path,
		[]byte(`server: address: ;`),
		0600,
	)

	require.NoError(t, err, "failed to write test config file %q", path)

	err = readFile(path, &initial)

	// ------ Проверки ------
	assert.Error(t, err, "reading an invalid file should return an error")
}

func TestReadMissingFile(t *testing.T) {
	initial := defaultRawConfig()

	err := readFile(filepath.Join(t.TempDir(), "non-existing.yaml"), &initial)

	// ------ Проверки ------
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
