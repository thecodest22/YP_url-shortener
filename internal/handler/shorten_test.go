package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodest22/YP_url-shortener/internal/service"
)

func TestHandlerShorten(t *testing.T) {
	// Запрос
	type req struct {
		contentType string
		body        string
	}
	// Задаваемые значения для возврата из Shorten
	type setShorten struct {
		retID  string
		retErr error
	}
	// Ожидаемое состояние Shorten после вызова
	type wantShorten struct {
		takeURL string
		calls   uint8
	}
	// Что ожидаем в респонсе от хэндлера
	type wantResp struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name        string
		req         req
		setShorten  setShorten
		wantShorten wantShorten
		wantResp    wantResp
	}{
		{
			name: "positive: valid content-type and simple URL",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://google.com",
			},
			setShorten: setShorten{
				retID:  "abc12345",
				retErr: nil,
			},
			wantShorten: wantShorten{
				takeURL: "https://google.com",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusCreated,
				response:    baseURL + "/abc12345",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive: boundary value for body size",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://ya.ru/" + strings.Repeat("a", bodyMaxBytes-len("https://ya.ru/")),
			},
			setShorten: setShorten{
				retID:  "abc12345",
				retErr: nil,
			},
			wantShorten: wantShorten{
				takeURL: "https://ya.ru/" + strings.Repeat("a", bodyMaxBytes-len("https://ya.ru/")),
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusCreated,
				response:    baseURL + "/abc12345",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive: body value whitespaces trimming",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        " \nhttps://ya.ru/\t\t", // До функции сервиса должно доехать без пробелов
			},
			setShorten: setShorten{
				retID:  "abc12345",
				retErr: nil,
			},
			wantShorten: wantShorten{
				takeURL: "https://ya.ru/",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusCreated,
				response:    baseURL + "/abc12345",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive: valid content-type and URL with extra path and params",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://google.com/test?foo=bar&spam=eggs",
			},
			setShorten: setShorten{
				retID:  "abc12345",
				retErr: nil,
			},
			wantShorten: wantShorten{
				takeURL: "https://google.com/test?foo=bar&spam=eggs",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusCreated,
				response:    baseURL + "/abc12345",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid body: empty body",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "",
			},
			setShorten:  setShorten{},  // До него дело не дойдет
			wantShorten: wantShorten{}, // calls должно быть 0
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid body: contains whitespaces only",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        " \n\t ", // Проверка еще до Shorten сервиса
			},
			setShorten:  setShorten{},  // До него дело не дойдет
			wantShorten: wantShorten{}, // calls должно быть 0
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid content-type: application/json",
			req: req{
				contentType: "application/json; charset=utf-8",
				body:        `{"url":"https://google.com"}`,
			},
			setShorten:  setShorten{},  // До него дело не дойдет
			wantShorten: wantShorten{}, // calls должно быть 0
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid headers: missing content-type",
			req: req{
				body: "https://google.com",
			},
			setShorten:  setShorten{},  // До него дело не дойдет
			wantShorten: wantShorten{}, // calls должно быть 0
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid body: content too large",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        strings.Repeat("a", bodyMaxBytes+1),
			},
			setShorten:  setShorten{},  // До него дело не дойдет
			wantShorten: wantShorten{}, // calls должно быть 0
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid body: not a URL passed",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "foo",
			},
			setShorten: setShorten{
				retID:  "",
				retErr: service.ErrInvalidURL,
			},
			wantShorten: wantShorten{
				takeURL: "foo",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "invalid url\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: invalid body: not a URL passed, wrapped error",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "foo",
			},
			setShorten: setShorten{
				retID:  "",
				retErr: fmt.Errorf("bad url: %w", service.ErrInvalidURL),
			},
			wantShorten: wantShorten{
				takeURL: "foo",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusBadRequest,
				response:    "invalid url\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: service err: ID generation attempts exhausted",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://ya.ru",
			},
			setShorten: setShorten{
				retID:  "",
				retErr: service.ErrIDGenerationExhausted,
			},
			wantShorten: wantShorten{
				takeURL: "https://ya.ru",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusServiceUnavailable,
				response:    "service unavailable\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: service err: ID generation attempts exhausted, wrapped error",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://ya.ru",
			},
			setShorten: setShorten{
				retID:  "",
				retErr: fmt.Errorf("ID generation failed: %w", service.ErrIDGenerationExhausted),
			},
			wantShorten: wantShorten{
				takeURL: "https://ya.ru",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusServiceUnavailable,
				response:    "service unavailable\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative: service err: unexpected error occurred",
			req: req{
				contentType: "text/plain; charset=utf-8",
				body:        "https://ya.ru",
			},
			setShorten: setShorten{
				retID:  "",
				retErr: errors.New("unexpected error"),
			},
			wantShorten: wantShorten{
				takeURL: "https://ya.ru",
				calls:   1,
			},
			wantResp: wantResp{
				code:        http.StatusInternalServerError,
				response:    "internal error\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for i, test := range tests {
		t.Run(getSubtestTitle(i+1, test.name), func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.req.body))
			if test.req.contentType != "" {
				request.Header.Set("Content-Type", test.req.contentType)
			}

			// Для отработки функции "shorten" у хэдлера должен быть сервис с соответствующей функцией.
			// Настраиваем мок сервиса, чтобы его функция Shorten отдавала нужные значения.
			svc := &mockURLShortener{
				sht: mockShorten{
					retID:  test.setShorten.retID,
					retErr: test.setShorten.retErr,
				},
			}

			router := NewRouter(logger, svc, baseURL)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			response := w.Result()
			defer response.Body.Close()

			// ------ Проверки -------
			assert.Equal(t, test.wantResp.code, response.StatusCode, "unexpected status code")
			assert.Equal(t, test.wantResp.contentType, response.Header.Get("Content-Type"), "unexpected content type")
			resBody, err := io.ReadAll(response.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.wantResp.response, string(resBody), "unexpected response body")
			// Дополнительно проверяем, что Shorten сервиса действительно вызвался 1 раз и с нужными аргументами
			assert.Equal(t, test.wantShorten.calls, svc.sht.calls, "unexpected calls")
			assert.Equal(t, test.wantShorten.takeURL, svc.sht.argURL, "unexpected URL argument")
		})
	}
}
