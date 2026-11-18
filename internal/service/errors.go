package service

import "errors"

var (
	ErrIDGenerationExhausted = errors.New("no free ID after retries")
	ErrInvalidURL            = errors.New("invalid URL value")
	ErrURLNotFound           = errors.New("entry with specified ID not found")
)
