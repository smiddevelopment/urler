package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/logger"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/config"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func init() {
	config.SetConfig()
	logger.InitLog()
}

func main() {

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(config.NetAddress.ServAddr, r)
	if err != nil {
		panic(err)
	}
}
