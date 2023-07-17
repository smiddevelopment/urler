package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

var port, resUrl *string

func init() {
	port = flag.String("a", "localhost:8888", "-a server port")
	resUrl = flag.String("b", "http://localhost:8080/", "-b result URL address")
	flag.Parse()
}

func main() {
	r := chi.NewRouter()

	decodeHandler := handler.New(resUrl)
	r.Post("/", decodeHandler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(*port, r)
	if err != nil {
		panic(err)
	}
}
