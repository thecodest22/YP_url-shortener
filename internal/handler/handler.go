package handler

import "net/http"

// NewRouter - создает корневой хэндлер с подключенными роутами.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	return mux
}
