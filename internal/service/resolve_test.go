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

func TestServiceResolve(t *testing.T) {
	readErr := errors.New("storage read unexpected error")
	tests := []struct {
		name           string
		id             string
		readRetURL     string
		readRetErr     error
		wantResolveURL string
		wantResolveErr error
		wantReadCalls  int
	}{
		{
			name:           "positive: URL found",
			id:             "aB12cd34",
			readRetURL:     "https://example.com/path?foo=bar",
			readRetErr:     nil,
			wantResolveURL: "https://example.com/path?foo=bar",
			wantResolveErr: nil,
			wantReadCalls:  1,
		},
		{
			name:           "negative: obviously invalid ID",
			id:             "invalid-ID",
			readRetURL:     "", // До Read хранилища дело не дойдет
			readRetErr:     nil,
			wantResolveURL: "",
			wantResolveErr: ErrURLNotFound,
			wantReadCalls:  0,
		},
		{
			name:           "negative: storage: key not found",
			id:             "aB12cd34",
			readRetURL:     "",
			readRetErr:     storage.ErrNotExists,
			wantResolveURL: "",
			wantResolveErr: ErrURLNotFound,
			wantReadCalls:  1,
		},
		{
			name:           "negative: storage: key not found with wrapped error",
			id:             "aB12cd34",
			readRetURL:     "",
			readRetErr:     fmt.Errorf("storage read err: %w", storage.ErrNotExists),
			wantResolveURL: "",
			wantResolveErr: ErrURLNotFound,
			wantReadCalls:  1,
		},
		{
			name:           "negative: storage: unexpected error",
			id:             "aB12cd34",
			readRetURL:     "",
			readRetErr:     readErr,
			wantResolveURL: "",
			wantResolveErr: readErr,
			wantReadCalls:  1,
		},
	}

	ctx := context.Background()

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			// Для отработки функции "Shorten" у сервиса должно быть хранилище с методами "Write" и "Read".
			// Настраиваем мок хранилища, чтобы его функция Read работала нужным образом.
			st := &mockURLStorage{
				readRetURL: test.readRetURL,
				readRetErr: test.readRetErr,
			}
			svc := NewURLShortener(st)

			url, err := svc.Resolve(ctx, test.id)
			calls := len(st.readCalledArgs)

			// ------ Проверки -------
			if test.wantResolveErr == nil {
				require.NoError(t, err, "must return no error")
				assert.Equal(t, test.wantResolveURL, url, "unexpected URL value returned")
			} else {
				require.ErrorIs(t, err, test.wantResolveErr, "unexpected error returned")
				assert.Empty(t, url, "URL must be empty on failure")
			}

			// Проверяем, что Read хранилища был вызван или не вызван, а если вызван, то с каким аргументом.
			assert.Equal(t, test.wantReadCalls, calls, "unexpected storage Read calls")
			if calls > 0 {
				assert.Equal(t, test.id, st.readCalledArgs[calls-1].key, "Read called with an unexpected key")
			}

			assert.Empty(t, st.writeCalledArgs, "Resolve must not call Write")
		})
	}
}
