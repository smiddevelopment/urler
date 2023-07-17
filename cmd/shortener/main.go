package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func init() {
	flag.String("a", "localhost:8080", "-a server address")
	flag.String("b", "http://localhost:8080", "-b result URL address")
}

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	servURL := flag.Lookup("a").Value.(flag.Getter).String()
	servURL = "localhost:8080"
	err := http.ListenAndServe(servURL, r)
	if err != nil {
		panic(err)
	}
}
