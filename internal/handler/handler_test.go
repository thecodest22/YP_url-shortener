// Переиспользуемые сущности для тестов данного пакета.
package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

const baseURL = "http://test"

// Структуры для описания работы методов сервиса:
// фиксация принятых аргументов и количества вызовов, настройка возвращаемых данных, и т.д.
// ---------------------
type mockShorten struct {
	argURL string
	retID  string
	retErr error
	calls  uint8
}

type mockResolve struct {
	argID  string
	retURL string
	retErr error
	calls  uint8
}

// ---------------------

// mockURLShortener. Мок для сервиса, чтобы подменить вызовы в хэндлере.
type mockURLShortener struct {
	sht mockShorten
	res mockResolve
}

// Объявляем методы на моке сервиса, чтобы он имплементировал ожидаемый интерфейс,
// возвращал какие-то преднастроенные данные, но реально ничего не делал (походы в хранилище и т.д.).
// ---------------------
func (s *mockURLShortener) Shorten(_ context.Context, url string) (string, error) {
	s.sht.calls++
	s.sht.argURL = url
	return s.sht.retID, s.sht.retErr
}

func (s *mockURLShortener) Resolve(_ context.Context, id string) (string, error) {
	s.res.calls++
	s.res.argID = id
	return s.res.retURL, s.res.retErr
}

func TestRouter_Routing(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "positive: valid path and method",
			method:     http.MethodGet,
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "positive: HEAD request to healthz",
			method:     http.MethodHead,
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "negative: not existing path",
			method:     http.MethodGet,
			path:       "/not-existing/path",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "negative: not existing path",
			method:     http.MethodPost,
			path:       "/not-existing/path",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "negative: not allowed method",
			method:     http.MethodPut,
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			svc := &mockURLShortener{}
			mux := NewRouter(logger, svc, baseURL)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest(test.method, baseURL+test.path, nil))

			response := w.Result()
			defer response.Body.Close()

			// ------- Проверки -------
			assert.Equal(t, test.wantStatus, response.StatusCode, "unexpected status code")
			assert.Zero(t, svc.sht.calls, "Shorten must not be called")
			assert.Zero(t, svc.res.calls, "Resolve must not be called")
		})
	}
}
