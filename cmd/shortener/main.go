package main

import (
	"net/http"

	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func main() {
	mux := http.NewServeMux() // Инициализация объекта сервера

	mux.HandleFunc(`/`, handler.RouteURL)     // Базовый маршрут
	mux.HandleFunc(`/{id}`, handler.RouteURL) // Маршрут с параметром ссылки

	err := http.ListenAndServe(`:8080`, mux) // Запуск сервера
	if err != nil {
		panic(err)
	}
}
