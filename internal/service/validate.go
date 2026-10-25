package service

import (
	"fmt"
	"net/url"
)

// isSupportedScheme - является ли схема URL поддерживаемой.
func isSupportedScheme(scheme string) bool {
	switch scheme {
	case "http", "https":
		return true
	default:
		return false
	}
}

// validateURL - валидирует переданный URL.
func validateURL(val string) error {
	u, err := url.ParseRequestURI(val)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	if !isSupportedScheme(u.Scheme) {
		return fmt.Errorf("%w: unsupported scheme %q", ErrInvalidURL, u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("%w: missing host", ErrInvalidURL)
	}

	return nil
}
