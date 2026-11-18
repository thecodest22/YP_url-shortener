package memory

import (
	"context"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodest22/YP_url-shortener/internal/storage"
	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestStorage_Write(t *testing.T) {
	tests := []struct {
		name     string
		initURLs map[string]string
		key      string
		url      string
		wantErr  error
		wantURLs map[string]string
	}{
		{
			name:     "positive: add key to empty map",
			initURLs: make(map[string]string), // Через make получаем пустую мапу, а не nil
			key:      "aB12cd34",
			url:      "https://example.com/path?foo=bar",
			wantErr:  nil,
			wantURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
			},
		},
		{
			name: "positive: add unique key to non-empty map",
			initURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
			},
			key:     "abcd1234",
			url:     "https://ya.ru",
			wantErr: nil,
			wantURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
				"abcd1234": "https://ya.ru",
			},
		},
		{
			name: "negative: add key that already exists",
			initURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
			},
			key:     "aB12cd34",
			url:     "https://ya.ru",
			wantErr: storage.ErrAlreadyExists,
			wantURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
			},
		},
	}
	ctx := context.Background()

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			st := New()
			st.urls = maps.Clone(test.initURLs)
			err := st.Write(ctx, test.key, test.url)
			assert.ErrorIs(t, err, test.wantErr, "Write unexpected error")
			assert.Equal(t, test.wantURLs, st.urls, "unexpected storage state after Write")
		})
	}
}

func TestStorage_Read(t *testing.T) {
	tests := []struct {
		name     string
		initURLs map[string]string
		key      string
		wantURL  string
		wantErr  error
	}{
		{
			name: "positive: read existing key",
			initURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
				"abcd1234": "https://ya.ru",
			},
			key:     "aB12cd34",
			wantURL: "https://example.com/path?foo=bar",
			wantErr: nil,
		},
		{
			name: "negative: read non-existing key",
			initURLs: map[string]string{
				"aB12cd34": "https://example.com/path?foo=bar",
				"abcd1234": "https://ya.ru",
			},
			key:     "fooBAR12",
			wantURL: "",
			wantErr: storage.ErrNotExists,
		},
	}
	ctx := context.Background()

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			st := New()
			st.urls = maps.Clone(test.initURLs)
			url, err := st.Read(ctx, test.key)
			assert.ErrorIs(t, err, test.wantErr, "unexpected Read error")
			assert.Equal(t, test.wantURL, url, "read unexpected data from storage")
			assert.Equal(t, test.initURLs, st.urls, "unexpected storage state after Read")
		})
	}
}
