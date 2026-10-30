package handler

import "net/http"

func (h *handler) health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		h.logger.Error("healthz write failed", "err", err)
	}
}
