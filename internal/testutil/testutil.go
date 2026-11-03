package testutil

import (
	"strconv"
	"strings"
)

// GetSubtestTitle Хэлпер для получения имени теста.
func GetSubtestTitle(num int, name string) string {
	var sb strings.Builder

	n := strconv.Itoa(num)
	// Оптимизация: заранее выделяем память под строку с учетом длин номера и имени.
	// 6 байт на ASCII-символы "test #" + пробел + длина номера + длина имени теста
	sb.Grow(7 + len(n) + len(name))

	sb.WriteString("test #")
	sb.WriteString(n)
	sb.WriteString(" ")
	sb.WriteString(name)

	return sb.String()
}
