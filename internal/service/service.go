package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/thecodest22/YP_url-shortener/internal/storage"
)

// FIXME Убрать это, когда будем тянуть конфиги извне
const (
	idLen      = 8  // Длина ключа, возвращаемого на фронт
	retriesMax = 10 // Количество попыток записи в хранилище
)

type urlStorage interface {
	Read(ctx context.Context, key string) (string, error)
	Write(ctx context.Context, key, val string) error
}

// URLShortener - сервис для сокращения URL.
type URLShortener struct {
	storage urlStorage
}

// NewURLShortener - создает структуру сервиса.
func NewURLShortener(storage urlStorage) *URLShortener {
	return &URLShortener{storage: storage}
}

// Shorten - выполняет сокращение URL, возвращает ID записи для последующего резолва.
// ErrIDGenerationExhausted - после всех попыток не удалось сгенерировать уникальный ID для хранения.
// ErrInvalidURL - передан невалидный URL.
func (s *URLShortener) Shorten(ctx context.Context, original string) (string, error) {
	if err := validateURL(original); err != nil {
		return "", err
	}

	for i := 0; i < retriesMax; i++ {
		idGen := generateID(idLen)

		err := s.storage.Write(ctx, idGen, original)
		if err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				continue
			}
			return "", fmt.Errorf("service url shorten failed: %w", err)
		}
		return idGen, nil
	}

	return "", ErrIDGenerationExhausted
}

// Resolve - получить значение оригинального URL по ID.
// ErrURLNotFound - если не нашли записи с таким ID.
func (s *URLShortener) Resolve(ctx context.Context, id string) (string, error) {
	res, err := s.storage.Read(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return "", ErrURLNotFound
		}
		return "", fmt.Errorf("service url resolve failed: %w", err)
	}

	return res, nil
}
