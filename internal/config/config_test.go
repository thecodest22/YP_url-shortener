package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecodest22/YP_url-shortener/internal/handler"
	"github.com/thecodest22/YP_url-shortener/internal/httpserver"
	"github.com/thecodest22/YP_url-shortener/internal/logger"
	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

// writeCfgFile - хэлпер, создает файл конфигурации.
func writeCfgFile(t *testing.T, name string, content []byte) (string, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	return path, os.WriteFile(path, content, 0600)
}

// TestLoad тест флоу сборки конфига
func TestLoad(t *testing.T) {
	defaults := defaultRawConfig().toConfig()

	fromFile := defaults
	fromFile.Handler.BaseURL = "https://file.com"

	fromCLI := defaults
	fromCLI.Handler.BaseURL = "https://cli.com"

	fromEnv := defaults
	fromEnv.Handler.BaseURL = "https://env.com"

	tests := []struct {
		name        string
		yamlContent []byte
		args        []string
		envs        map[string]string
		want        Config
	}{
		{
			name:        "positive: external values absence preserve program defaults",
			yamlContent: nil,
			args:        []string{},
			envs:        map[string]string{},
			want:        defaults,
		},
		{
			name: "positive: file values override program defaults",
			yamlContent: []byte(`
                handler:
                  base_url: "https://file.com"
            `),
			args: []string{},
			envs: map[string]string{},
			want: fromFile,
		},
		{
			name:        "positive: CLI values override program defaults",
			yamlContent: nil,
			args:        []string{"-b", "https://cli.com"},
			envs:        map[string]string{},
			want:        fromCLI,
		},
		{
			name:        "positive: env values override program defaults",
			yamlContent: nil,
			args:        []string{},
			envs: map[string]string{
				"BASE_URL": "https://env.com",
			},
			want: fromEnv,
		},
		{
			name: "positive: CLI values override file values",
			yamlContent: []byte(`
                handler:
                  base_url: "https://file.com"
            `),
			args: []string{"-b", "https://cli.com"},
			envs: map[string]string{},
			want: fromCLI,
		},
		{
			name:        "positive: env values override CLI values",
			yamlContent: nil,
			args:        []string{"-b", "https://cli.com"},
			envs: map[string]string{
				"BASE_URL": "https://env.com",
			},
			want: fromEnv,
		},
		{
			name: "positive: env overrides file and CLI",
			yamlContent: []byte(`
                handler:
                  base_url: "https://file.com"
            `),
			args: []string{"-b", "https://cli.com"},
			envs: map[string]string{
				"BASE_URL": "https://env.com",
			},
			want: fromEnv,
		},
		{
			name: "positive: single vars from all sources",
			yamlContent: []byte(`
                server:
                  address: "localhost:8888"
            `),
			args: []string{"--log-level", "ERROR"},
			envs: map[string]string{
				"BASE_URL": "https://env.com",
			},
			want: Config{
				Server: httpserver.Config{
					Addr:              "localhost:8888",
					ReadHeaderTimeout: defaultReadHeaderTimeout,
					ShutdownTimeout:   defaultShutdownTimeout,
				},
				Logger: logger.Config{
					Level: "ERROR",
				},
				Handler: handler.Config{
					BaseURL: "https://env.com",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d %s", i+1, test.name), func(t *testing.T) {
			stashEnvs(
				t,
				[]string{
					"LOG_LEVEL",
					"SERVER_ADDRESS",
					"BASE_URL",
				},
			)

			if test.yamlContent != nil {
				f := "config.yaml"
				path, err := writeCfgFile(t, f, test.yamlContent)
				require.NoError(t, err, "failed to write test config file %q", f)
				test.args = append(test.args, "-c", path)
			}

			if test.envs != nil {
				setEnvs(t, test.envs)
			}

			gotCfg, err := Load(test.args)

			require.NoError(t, err, "unexpected error")
			assert.Equal(t, test.want, gotCfg, "got incorrect configuration")
		})
	}
}

// TestLoadErrors проверяет возврат ошибки и пустой структуры Config из Load.
func TestLoadErrors(t *testing.T) {
	stashEnvs(t, []string{
		"SERVER_ADDRESS",
		"BASE_URL",
		"LOG_LEVEL",
	})

	missingPath := filepath.Join(t.TempDir(), "missing.yaml")

	tests := []struct {
		name    string
		args    []string
		wantErr testutil.ErrorCheck
	}{
		{
			name:    "unknown flag",
			args:    []string{"--unknown-flag"},
			wantErr: testutil.ErrorAs[*pflag.NotExistError],
		},
		{
			name:    "missing config file",
			args:    []string{"--config", missingPath},
			wantErr: testutil.ErrorIs(fs.ErrNotExist),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d %s", i+1, test.name), func(t *testing.T) {
			got, err := Load(test.args)

			test.wantErr(t, err)
			assert.Equal(t, Config{}, got)
		})
	}
}
