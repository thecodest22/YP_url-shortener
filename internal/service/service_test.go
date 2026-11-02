// Переиспользуемые сущности для тестов данного пакета.
package service

import "context"

// Структуры для описания работы методов хранилища сервиса:
// фиксация принятых аргументов и количества вызовов, настройка возвращаемых данных, и т.д.
// ---------------------
type writeCallArgs struct {
	key string
	val string
}
type readCallArgs struct {
	key string
}

type mockURLStorage struct {
	// Сюда будем при инициализации явно добавлять ошибки или nil, соответствующие по порядку вызовам.
	writeRetErrs []error
	readRetErrs  []error
	// Сюда будут записываться аргументы функции на каждый "вызов".
	writeCalledArgs []writeCallArgs
	readCalledArgs  []readCallArgs
}

// Объявляем методы на моке хранилища, чтобы оно имплементировало ожидаемый интерфейс,
// возвращало какие-то преднастроенные данные, но реально не делало запросы в хранилище.
// ---------------------
func (s *mockURLStorage) Write(_ context.Context, key, val string) error {
	callIndex := len(s.writeCalledArgs)
	s.writeCalledArgs = append(s.writeCalledArgs, writeCallArgs{
		key: key,
		val: val,
	})

	if callIndex < len(s.writeRetErrs) {
		return s.writeRetErrs[callIndex]
	}
	return nil
}

func (s *mockURLStorage) Read(_ context.Context, _ string) (string, error) {
	panic("unexpected Read call in Shorten test")
}
