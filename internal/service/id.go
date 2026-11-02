package service

import (
	"math/rand/v2"
	"strings"
)

// idCharset - символы для ключей: латиница в обоих регистрах и цифры (всего 62 символа).
// FIXME убрать это, когда будем тянуть конфиги извне
const idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateID - сгенерировать ID-строку заданной длины.
func generateID(length int) string {
	var sb strings.Builder
	sb.Grow(length) // Оптимизируем выделение памяти под билдер.

	for i := 0; i < length; i++ {
		randIdx := rand.IntN(len(idCharset)) // Генерит в интервале [0, X)
		sb.WriteByte(idCharset[randIdx])     // Пишем чистые байты в буфер: так быстрее и меньше памяти
	}

	return sb.String()
}
