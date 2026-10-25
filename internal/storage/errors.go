package storage

import "errors"

var (
	ErrNotExists     = errors.New("key not exists")     // Нет такого ключа в хранилище
	ErrAlreadyExists = errors.New("key already exists") // На случай совпадения сгенерированного ключа
)
