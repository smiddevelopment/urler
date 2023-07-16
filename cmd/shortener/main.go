package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.EncodeURL)
	mux.HandleFunc(`/{id}`, handler.DecodeURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
