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

func TestIsIDValid(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "positive: mixed letters and digits",
			id:   "abcXYZ09",
			want: true,
		},
		{
			name: "positive: Latin letters in lower case only",
			id:   "abcdefgh",
			want: true,
		},
		{
			name: "positive: Latin letters in upper case only",
			id:   "ABCDEFGH",
			want: true,
		},
		{
			name: "positive: digits only",
			id:   "01234567",
			want: true,
		},
		{
			name: "negative: empty value",
			id:   "",
			want: false,
		},
		{
			name: "negative: unsupported char at start",
			id:   "-abcDE12",
			want: false,
		},
		{
			name: "negative: unsupported char in middle",
			id:   "abc-DE12",
			want: false,
		},
		{
			name: "negative: unsupported char at end",
			id:   "abcDE12-",
			want: false,
		},
		{
			name: "negative: whitespace at start",
			id:   " abcDE12",
			want: false,
		},
		{
			name: "negative: whitespace in middle",
			id:   "abc DE12",
			want: false,
		},
		{
			name: "negative: non-latin chars",
			id:   "123абвABC",
			want: false,
		},
		{
			name: "negative: non-latin chars",
			id:   "123💩🤡abc",
			want: false,
		},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, isIDValid(test.id))
		})
	}
}
