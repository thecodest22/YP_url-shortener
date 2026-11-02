package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecodest22/YP_url-shortener/internal/storage"
	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestServiceShorten(t *testing.T) {
	writeErr := errors.New("storage write unexpected error")
	tests := []struct {
		name           string
		originalURL    string
		writeErrs      []error
		wantWriteCalls int
		wantShortenErr error
	}{
		{
			name:           "positive: write success on the first call",
			originalURL:    "https://google.com",
			writeErrs:      []error{nil},
			wantWriteCalls: 1,
			wantShortenErr: nil,
		},
		{
			name:           "positive: write success on the second call",
			originalURL:    "https://ya.ru",
			writeErrs:      []error{storage.ErrAlreadyExists, nil},
			wantWriteCalls: 2,
			wantShortenErr: nil,
		},
		{
			name:        "positive: write success after one wrapped error",
			originalURL: "https://ya.ru",
			writeErrs: []error{
				fmt.Errorf("write err: %w", storage.ErrAlreadyExists),
				nil,
			},
			wantWriteCalls: 2,
			wantShortenErr: nil,
		},
		{
			name:        "positive: write success with boundary retries",
			originalURL: "https://ya.ru",
			writeErrs: []error{
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				nil,
			},
			wantWriteCalls: 10,
			wantShortenErr: nil,
		},
		{
			name:        "negative: write failed due to retries exceeded",
			originalURL: "https://google.com",
			writeErrs: []error{
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
			},
			wantWriteCalls: 10,
			wantShortenErr: ErrIDGenerationExhausted,
		},
		{
			name:        "negative: write failed due to unexpected error",
			originalURL: "http://foo.org",
			writeErrs: []error{
				storage.ErrAlreadyExists,
				storage.ErrAlreadyExists,
				writeErr,
			},
			wantWriteCalls: 3,
			wantShortenErr: writeErr,
		},
		{
			name:           "negative: invalid original URL",
			originalURL:    "foo",
			writeErrs:      nil, // При невалидном URL до работы с хранилищем не дойдет
			wantWriteCalls: 0,
			wantShortenErr: ErrInvalidURL,
		},
	}

	ctx := context.Background()

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			// Для отработки функции "Shorten" у сервиса должно быть хранилище с методами "Write" и "Read".
			// Настраиваем мок хранилища, чтобы его функция Write работала нужным образом.
			st := &mockURLStorage{
				writeRetErrs: test.writeErrs,
			}
			svc := NewURLShortener(st)

			id, err := svc.Shorten(ctx, test.originalURL)
			calls := len(st.writeCalledArgs)

			// ------ Проверки -------
			if test.wantShortenErr == nil {
				require.NoError(t, err, "must return no error")
				// Если Shorten отработал без ошибок, делаем проверки результата:
				// 1) сравниваем длину полученного ID
				assert.Len(t, id, idLen, "unexpected length of returned ID")
				// 2) если вообще были вызовы Write, проверяем, что по итогу Shorten вернула тот же ID,
				// с которым в конечном итоге отработала Write.
				if calls > 0 {
					assert.Equal(t, st.writeCalledArgs[calls-1].key, id, "ID from Shorten must match ID into Write")
				}
			} else {
				require.ErrorIs(t, err, test.wantShortenErr, "unexpected error return from Shorten")
				assert.Empty(t, id, "ID must be empty on failure")
			}
			// Здесь проверяем через require, чтобы сразу тест отвалился, иначе дальше
			// рискуем получить панику из-за несуществующего индекса.
			require.Equal(t, test.wantWriteCalls, calls, "unexpected storage Write calls")

			for _, call := range st.writeCalledArgs {
				assert.Equal(t, test.originalURL, call.val, "storage Write called with unexpected argument")
			}

			assert.Empty(t, st.readCalledArgs, "Shorten must not call Read")
		})
	}
}
