package handler

import (
	"errors"
	"fmt"
	"github.com/thecodest22/YP_url-shortener/internal/service"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

// bodyMaxBytes - максимальный размер тела запроса. Ставлю мегабайт от балды.
const bodyMaxBytes = 1024 * 1024

// isPlainText - проверяет, является ли Content-Type "text/plain".
func isPlainText(contentType string) bool {
	// Могут прислать не просто значение, например, "text/plain", а с доп. параметрами,
	// например, "text/plain; charset=utf-8". Поэтому заюзаем либу, которая
	// распарсит значение и отделит параметры от типа. Проверка будет регистронезависимой.
	mediaType, _, err := mime.ParseMediaType(contentType)
	return err == nil && mediaType == "text/plain"
}

// parseShortenRequest - достает сокращаемый URL из тела запроса.
func parseShortenRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	// Проверяем заголовок "Content-Type"
	if ct := r.Header.Get("Content-Type"); !isPlainText(ct) {
		// Тут нет типа ошибки, поэтому врапать нечего.
		// Если захотим снаружи отличать ошибку и возвращать 415, сделаем кастомную ошибку.
		return "", fmt.Errorf("unsupported content type %q", ct)
	}

	// Навешиваем ограничитель размера чтения тела для защиты.
	// Тогда при прочтении тела с превышенным размером получим *http.MaxBytesError.
	r.Body = http.MaxBytesReader(w, r.Body, bodyMaxBytes)

	// NOTE! Тут в байтах отдает, поэтому нужно будет привести к стринге!
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Ошибка есть, поэтому врапаем ее. Сейчас по ТЗ на любой кривой запрос отдаем 400,
		// но если захотим отловить превышение размера, то достанем через errors.As и отдадим 413.
		return "", fmt.Errorf("read body: %w", err)
	}

	payload := strings.TrimSpace(string(body))
	if payload == "" {
		return "", errors.New("empty body")
	}

	return payload, nil
}

// shorten - обрабатывает запрос на укорачивание URL.
func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := parseShortenRequest(w, r)
	if err != nil {
		// По ТЗ любая кривой запрос - отдаем 400.
		h.logger.Debug("bad shorten request", "err", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Если все ОК, запускаем функцию сервиса.
	id, err := h.svc.Shorten(r.Context(), originalURL)
	if err != nil {
		// Здесь могут быть разные ошибки, поэтому сравниваем и отдаем разные коды.
		switch {
		// Кривой URL - это клиентский косяк, отдаем 400.
		// Логируем уровнем DEBUG, т.к. штатная ситуация.
		case errors.Is(err, service.ErrInvalidURL):
			h.logger.Debug("invalid url", "err", err)
			http.Error(w, "invalid url", http.StatusBadRequest)

		// Не удалось сгенерить уникальный ID - наш косяк, но можно снова попробовать, поэтому 503.
		// Логируем уровнем ERROR, т.к. такое надо в логах видеть.
		case errors.Is(err, service.ErrIDGenerationExhausted):
			h.logger.Error("shorten url", "err", err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)

		// Неизвестный сбой - все плохо, клиенту 500, в логи ERROR.
		default:
			h.logger.Error("shorten url", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	// Соединяем базовый URL со сгенерированным ID, чтобы получить строку, как в ТЗ.
	// Заюзаем либу, т.к. она корректно соединит через слэши. На ошибку отдаем 500.
	shortURL, err := url.JoinPath(h.baseURL, id)
	if err != nil {
		h.logger.Error("build short url", "base_url", h.baseURL, "id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Задаем заголовки, отправляем заголовки и ответ.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	// Можно писать через io.WriteString. Она проверяет у писателя метод, позволяющий писать без перевода в байты.
	// Тогда, в отличие от w.Write([]byte(string)), не будет создаваться копия строки.
	// Если же такого метода нет, то она вызовет штатный w.Write([]byte(string)).
	if _, err := io.WriteString(w, shortURL); err != nil {
		// Ошибка при записи - это, скорее всего, уже отключился клиент.
		// Ответ изменить нельзя, можем только залогировать.
		h.logger.Debug("write shorten response", "err", err)
	}
}
