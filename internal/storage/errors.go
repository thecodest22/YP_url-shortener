package storage

import "errors"

var (
	ErrNotFound      = errors.New("value not found")    // Значение не найдено по ключу
	ErrAlreadyExists = errors.New("key already exists") // На случай совпадения сгенерированного ключа
)
