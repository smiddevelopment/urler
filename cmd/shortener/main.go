package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

func main() {
	e := echo.New()
	e.POST(`/`, handler.EncodeUrl)
	e.GET(`/{id}`, handler.DecodeUrl)

	e.Logger.Fatal(e.Start(":8080"))
}
