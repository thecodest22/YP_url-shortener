package handler

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")
	ErrMissingURLHostname   = errors.New("missing URL hostname")
)

// validateBaseURL валидирует базовый URL для формирования ссылок.
func validateBaseURL(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}

	// Тут врапнем, чтобы в сообщении содержалось, какой протокол был на входе
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: %q", ErrUnsupportedURLScheme, u.Scheme)
	}

	if u.Hostname() == "" {
		return ErrMissingURLHostname
	}

	return nil
}
