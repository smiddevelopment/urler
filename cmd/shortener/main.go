package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/gzipMiddleware"

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
	r.Use(gzipMiddleware.GzipMiddleware)
	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)
	r.Post("/api/shorten", handler.EncodeURLJSON)

	err := http.ListenAndServe(config.NetAddress.ServAddr, r)
	if err != nil {
		panic(err)
	}
}
