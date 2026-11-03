package handler

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test",
			want: want{
				code:        http.StatusOK,
				response:    "OK",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	// Создаем логгер, т.к. он нужен для инициализации структуры хэндлера.
	// Он будет писать в никуда, т.к. io.Discard.Write устроена так, что она
	// просто возвращает успешный ответ, никак не используя и никуда не сохраняя переданные байты.
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем тестовый запрос из специальной либы. Он отдает тот же Request,
			// что и для боевого кода, со всеми полями и методами.
			request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

			// Создаем ResponseRecorder. Он реализует интерфейс http.ResponseWriter,
			// как в боевом коде.
			w := httptest.NewRecorder()

			// Функция "healthz" не работает с сервисом, поэтому передаем nil вместо сервиса.
			router := NewRouter(logger, nil, baseURL)

			// Вызовем не просто метод на структуре хэндлера, а именно через ServeHTTP,
			// поскольку так мы проверим еще и правильность маршрута, и допустимость метода.
			router.ServeHTTP(w, request)

			// Вызываем Result(), он отдаст нам *http.Response - как в боевом коде.
			// Из него будем читать, что вернул хэндлер и проверять результат.
			response := w.Result()
			defer response.Body.Close()

			// ------ Проверки -------
			assert.Equal(t, test.want.code, response.StatusCode, "unexpected status code")
			assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"), "unexpected content type")
			resBody, err := io.ReadAll(response.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody), "unexpected response body")
		})
	}
}
