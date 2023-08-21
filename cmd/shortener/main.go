package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/storage"

	"github.com/smiddevelopment/urler.git/internal/app/gzipmiddleware"

	"github.com/smiddevelopment/urler.git/internal/app/logger"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/config"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func init() {
	config.SetConfig()
	logger.InitLog()
	storage.InitStore()
}

func main() {

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(gzipmiddleware.GzipMiddleware)
	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)
	r.Post("/api/shorten", handler.EncodeURLJSON)
	r.Get("/ping", handler.PingDB)

	err := http.ListenAndServe(config.ServerConfig.ServAddr, r)
	if err != nil {
		panic(err)
	}
}
