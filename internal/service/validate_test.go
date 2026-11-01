package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "positive: valid HTTP URL",
			url:     "http://ya.ru",
			wantErr: nil,
		},
		{
			name:    "positive: valid HTTPS URL",
			url:     "https://google.com",
			wantErr: nil,
		},
		{
			name:    "positive: valid HTTPS URL with path and params",
			url:     "https://google.com/some-path?foo=bar&spam=eggs",
			wantErr: nil,
		},
		{
			name:    "negative: empty value",
			url:     "",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "negative: not a URL",
			url:     "foo",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "negative: unsupported scheme",
			url:     "ftp://example.com/file",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "negative: missing host",
			url:     "http://",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "negative: missing host",
			url:     "http:///some-path",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "negative: invalid percentage encoding",
			url:     "http://foo.bar/%zz",
			wantErr: ErrInvalidURL,
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			err := validateURL(test.url)

			if test.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.wantErr)
			}
		})
	}
}
