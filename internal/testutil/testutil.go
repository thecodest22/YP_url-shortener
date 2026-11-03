package testutil

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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

// ErrorCheck - кастомный тип-фукнция для более удобной работы с проверками ошибок.
type ErrorCheck func(t *testing.T, err error)

// ErrorAs - проверка ошибки по типу.
// [E error] - это дженерик (почти как в питоне). Тип проверяется во время компиляции.
// В нашем случае сюда может залететь кастомная структура-error,
// и тогда при проверке типов будет сравнение по интерфейсу.
// Если все ОК, сможем проверять переданный аргумент через ErrorAs.
// NOTE!!! Выбор типа еще не вызывает функцию, а только инициализирует! Например:
// check := ErrorAs[*pflag.NotExistError] - здесь только получаем функцию.
// check(t, err) - только затем вызываем с нужными аргументами.
func ErrorAs[E error](t *testing.T, err error) {
	// Это помечает текущую функцию как вспомогательную для тестов.
	// Когда пакет testing определяет, к какой строке привязать сообщение об ошибке,
	// он пропускает такие функции.
	t.Helper()

	var target E
	require.ErrorAs(t, err, &target)
}

// ErrorIs просто возвращает функцию проверки ошибок.
// Она удобна для использования в тестах и может быть присвоена полю в таблице тестов,
// когда надо ловить простые ошибки (не структуры).
func ErrorIs(target error) ErrorCheck {
	return func(t *testing.T, err error) {
		t.Helper()

		require.ErrorIs(t, err, target)
	}
}
