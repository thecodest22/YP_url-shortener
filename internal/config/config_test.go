package config

import (
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
			name: "positive: no args passed",
			args: []string{},
			wantCfg: Config{
				BaseURL:    "http://localhost:8080",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with full name",
			args: []string{"--base-url", "http://foo.com"},
			wantCfg: Config{
				BaseURL:    "http://foo.com",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed with shorthand",
			args: []string{"-b", "http://foo.com"},
			wantCfg: Config{
				BaseURL:    "http://foo.com",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				BaseURL:    "http://ya.ru",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"-b", "http://foo.com", "--base-url", "http://ya.ru"},
			wantCfg: Config{
				BaseURL:    "http://ya.ru",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: base URL passed twice",
			args: []string{"--base-url", "http://foo.com", "-b", "http://ya.ru"},
			wantCfg: Config{
				BaseURL:    "http://ya.ru",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "positive: ignorable position args passed",
			args: []string{"foo", "bar"},
			wantCfg: Config{
				BaseURL:    "http://localhost:8080",
				ServerAddr: "localhost:8080",
			},
			wantErr: nil,
		},
		{
			name: "negative: unknown args passed",
			args: []string{"--foo", "bar"},
			wantCfg: Config{
				BaseURL:    "",
				ServerAddr: "",
			},
			wantErr: testutil.ErrorAs[*pflag.NotExistError], // Тут функция только инициализируется из дженерика
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			// Т.к. мы работаем с переменными пакета,
			// нужно сначала сохранить состояние, которое было до теста:
			savedCfg, savedFS := cfg, fs

			// Регистрируем восстановление с помощью специальной функции.
			// Она выполнится после завершения теста (даже если упадет на require).
			t.Cleanup(func() {
				cfg = savedCfg
				fs = savedFS
			})

			// Приводим переменные пакета в дефолтное состояние для текущего теста.
			cfg = Config{}
			fs = newFlagSet(&cfg)

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
