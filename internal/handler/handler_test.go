// Переиспользуемые сущности для тестов данного пакета.
package handler

import (
	"context"
)

const baseURL = "http://test"

// Структуры для описания работы методов сервиса:
// фиксация принятых аргументов и количества вызовов, настройка возвращаемых данных, и т.д.
// ---------------------
type mockShorten struct {
	argURL string
	retID  string
	retErr error
	calls  uint8
}

type mockResolve struct {
	argID  string
	retURL string
	retErr error
	calls  uint8
}

// ---------------------

// mockURLShortener. Мок для сервиса, чтобы подменить вызовы в хэндлере.
type mockURLShortener struct {
	sht mockShorten
	res mockResolve
}

// Объявляем методы на моке сервиса, чтобы он имплементировал ожидаемый интерфейс,
// возвращал какие-то преднастроенные данные, но реально ничего не делал (походы в хранилище и т.д.).
// ---------------------
func (s *mockURLShortener) Shorten(_ context.Context, url string) (string, error) {
	s.sht.calls++
	s.sht.argURL = url
	return s.sht.retID, s.sht.retErr
}

func (s *mockURLShortener) Resolve(_ context.Context, id string) (string, error) {
	s.res.calls++
	s.res.argID = id
	return s.res.retURL, s.res.retErr
}
