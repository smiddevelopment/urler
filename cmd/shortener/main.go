package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func init() {
	flag.String("a", "localhost:8080", "-a server address")
	flag.String("b", "http://localhost:8080", "-b result URL address")
}

func main() {
	flag.Parse()

	servURL, exist := os.LookupEnv("SERVER_ADDRESS")
	if exist {
		fmt.Print(servURL)
	}
	if servURL == "" {
		servURL = flag.Lookup("a").Value.(flag.Getter).String()
	}

	r := chi.NewRouter()

	r.Post("/", handler.EncodeURL)
	r.Get("/{id}", handler.DecodeURL)

	err := http.ListenAndServe(servURL, r)
	if err != nil {
		panic(err)
	}
}
