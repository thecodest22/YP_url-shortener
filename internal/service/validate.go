package service

import (
	"fmt"
	"net/url"
	"strings"
)

// idCharsetRegex - скомпилированная регулярка для проверки допустимых символов в ID.
//var idCharsetRegex = regexp.MustCompile(
//	// regexp.QuoteMeta() экранирует символы, чтобы они не интерпретировались
//	// регуляркой как управляющие (вдруг добавим дефисы и прочую шелуху).
//	// Получится паттерн: ^[abc...789\-_]+$
//	"^[" + regexp.QuoteMeta(idCharset) + "]+$",
//)

// isSupportedScheme - является ли схема URL поддерживаемой.
func isSupportedScheme(scheme string) bool {
	switch scheme {
	case "http", "https":
		return true
	default:
		return false
	}
}

// isIDValid - валиден ли ID для резолва оригинального URL.
// Значение должно содержать только символы, из которых мы генерим ID при укорачивании.
func isIDValid(id string) bool {
	if id == "" {
		return false
	}

	// Проверка циклом (быстрее, но больше кода).
	// Когда идем циклом по строке, на выходе получаем руны.
	// Поэтому используем специальный оптимизированный поиск руны в строке.
	for _, c := range id {
		if !strings.ContainsRune(idCharset, c) {
			return false
		}
	}

	return true

	// Проверка регуляркой (медленнее, но лаконичнее).
	//return idCharsetRegex.MatchString(id)
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
