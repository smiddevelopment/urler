package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/config"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func init() {
	config.SetConfig()
}

func main() {
	r := chi.NewRouter()

	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(config.NetAddress.ServAddr, r)
	if err != nil {
		panic(err)
	}
}
