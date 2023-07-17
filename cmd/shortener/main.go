package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handler.RouteURL)
	//mux.HandleFunc(`/{id}`, handler.RouteURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
