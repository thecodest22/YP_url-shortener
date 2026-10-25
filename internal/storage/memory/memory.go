package memory

import (
	"fmt"
	"github.com/thecodest22/YP_url-shortener/internal/storage"
	"sync"
)

// Storage - хранилище сокращенных URL'ов в памяти.
type Storage struct {
	mu   sync.RWMutex // Инициализировать не нужно, zero value уже рабочий!
	urls map[string]string
}

// New - инициализирует хранилище сокращенных URL'ов в памяти.
func New() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

// Write - записывает значение по ключу.
// storage.ErrAlreadyExists - если в хранилище уже есть запись под таким ключом.
func (s *Storage) Write(key, val string) error {
	// Блочим на запись, иначе будут конфликты между потоками!
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, не существует ли уже запись с таким ключом.
	if _, exists := s.urls[key]; exists {
		return fmt.Errorf("storage write failed: %w", storage.ErrAlreadyExists)
	}

	// Если все ОК, пишем значение.
	s.urls[key] = val
	return nil
}

// Read - читает значение по ключу.
// Пустая строка и storage.ErrNotFound - если значение не найдено.
func (s *Storage) Read(key string) (string, error) {
	// Блочим на чтение. Читать параллельно можно, но будут конфликты,
	// если кто-то будет писать!
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, exists := s.urls[key]
	if !exists {
		return "", fmt.Errorf("storage read failed: %w", storage.ErrNotFound)
	}

	return val, nil
}
