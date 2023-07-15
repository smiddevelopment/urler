package main

import (
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

var port, resUrl *string

func init() {
	port = flag.String("a", "8080", "-a is the server path")
	resUrl = flag.String("b", "http://localhost:8080", "-b <base address result>")
	flag.Parse()
}

func main() {
	e := echo.New()

	decodeHandler := handler.New(resUrl)
	e.POST(`/`, decodeHandler.EncodeUrl)
	e.GET(`/{id}`, decodeHandler.DecodeUrl)

	defaultPort := "5000"

	if port != nil {
		defaultPort = *port
	}

	e.Logger.Fatal(e.Start(":" + defaultPort))
}
