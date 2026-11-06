package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thecodest22/YP_url-shortener/internal/service"
	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestHandlerResolve(t *testing.T) {
	// Запрос
	type req struct {
		path string
	}
	// Задаваемые значения для возврата из Resolve
	type setResolve struct {
		retURL string
		retErr error
	}
	// Ожидаемое состояние Resolve после вызова
	type wantResolve struct {
		takeID string
		calls  uint8
	}
	// Что ожидаем в респонсе от хэндлера
	type wantResp struct {
		code     int
		location string
	}

	tests := []struct {
		name        string
		req         req
		setResolve  setResolve
		wantResolve wantResolve
		wantResp    wantResp
	}{
		{
			name: "positive: valid ID passed",
			req: req{
				path: "/abc12345",
			},
			setResolve: setResolve{
				retURL: "https://ya.ru/test?foo=bar&spam=eggs",
				retErr: nil,
			},
			wantResolve: wantResolve{
				takeID: "abc12345",
				calls:  1,
			},
			wantResp: wantResp{
				code:     http.StatusTemporaryRedirect,
				location: "https://ya.ru/test?foo=bar&spam=eggs",
			},
		},
		{
			name: "negative: ID is valid but not found in a storage",
			req: req{
				path: "/abc12345",
			},
			setResolve: setResolve{
				retURL: "",
				retErr: service.ErrURLNotFound,
			},
			wantResolve: wantResolve{
				takeID: "abc12345",
				calls:  1,
			},
			wantResp: wantResp{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "negative: ID is valid but not found in a storage",
			req: req{
				path: "/abc12345",
			},
			setResolve: setResolve{
				retURL: "",
				// Для врапнутой ошибки тоже сделаем тесткейс
				retErr: fmt.Errorf("resolve failed: %w", service.ErrURLNotFound),
			},
			wantResolve: wantResolve{
				takeID: "abc12345",
				calls:  1,
			},
			wantResp: wantResp{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "negative: invalid ID: unsupported chars",
			req: req{
				path: "/abc-1234",
			},
			setResolve: setResolve{
				retURL: "",
				retErr: service.ErrURLNotFound,
			},
			wantResolve: wantResolve{
				takeID: "abc-1234",
				calls:  1,
			},
			wantResp: wantResp{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "negative: service error: unexpected error occurred",
			req: req{
				path: "/abc12345",
			},
			setResolve: setResolve{
				retURL: "",
				retErr: errors.New("unexpected error"),
			},
			wantResolve: wantResolve{
				takeID: "abc12345",
				calls:  1,
			},
			wantResp: wantResp{
				code:     http.StatusInternalServerError,
				location: "",
			},
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, baseURL+test.req.path, nil)

			// Для отработки функции "resolve" у хэдлера должен быть сервис с соответствующей функцией.
			// Настраиваем мок сервиса, чтобы его функция Resolve отдавала нужные значения.
			svc := &mockURLShortener{
				res: mockResolve{
					retURL: test.setResolve.retURL,
					retErr: test.setResolve.retErr,
				},
			}

			router := NewRouter(logger, svc, baseURL)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			response := w.Result()
			defer response.Body.Close()

			// ------ Проверки -------
			assert.Equal(t, test.wantResp.code, response.StatusCode, "unexpected status code")
			// Проверять отсутствие заголовка как response.Header.Get("Location") == "" будет некорректно,
			// т.к. и при отсутствии заголовка, и реально при пустой строке будет выглядеть одинаково.
			// У нас в таблице тестов "" - это именно отсутствие, поэтому проверяем на отсутствие ключа.
			if test.wantResp.location == "" {
				assert.NotContains(t, response.Header, "Location", "Location header must be absent")
			} else {
				assert.Equal(t, test.wantResp.location, response.Header.Get("Location"), "unexpected location header")
			}
			// Дополнительно проверяем, что Resolve сервиса действительно вызвался 1 раз и с нужными аргументами
			assert.Equal(t, test.wantResolve.calls, svc.res.calls, "unexpected calls")
			assert.Equal(t, test.wantResolve.takeID, svc.res.argID, "unexpected ID argument")
		})
	}
}
