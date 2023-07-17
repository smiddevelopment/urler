package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

var port *string

func init() {
	port = flag.String("a", "localhost:8888", "-a server port")
	flag.String("b", "http://localhost:8080/", "-b result URL address")
}

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(*port, r)
	if err != nil {
		panic(err)
	}
}
