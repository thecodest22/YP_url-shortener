package main

import (
	"log"

	"github.com/thecodest22/YP_url-shortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("service fatal error: %s", err)
	}
}
