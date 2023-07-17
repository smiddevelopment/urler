package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
