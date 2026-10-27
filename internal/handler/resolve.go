package handler

import (
	"errors"
	"net/http"

	"github.com/thecodest22/YP_url-shortener/internal/service"
)

// resolve - обрабатывает запрос на возврат оригинального URL по ID сокращенного.
func (h *handler) resolve(w http.ResponseWriter, r *http.Request) {
	// Думаю, тут не может быть каких-то внештатных ситуаций, т.к. если бы не было
	// передано значение в /{id}, маршрут бы не совпал и сюда бы дело не дошло.
	id := r.PathValue("id")

	resURL, err := h.svc.Resolve(r.Context(), id)
	if err != nil {
		// По ТЗ на любой кривой запрос отдаем 400.
		if errors.Is(err, service.ErrURLNotFound) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Если непредвиденная ошибка, 500
		h.logger.Error("resolve ID", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Возьмем хэлпер для редиректа из уроков.
	http.Redirect(w, r, resURL, http.StatusTemporaryRedirect)
}
