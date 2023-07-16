package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/smiddevelopment/urler.git/internal/app/handler"
)

var options struct {
	port   string
	resUrl string
}

func init() {
	flag.StringVar(&options.port, "a", "8080", "-a is the server path")
	flag.StringVar(&options.resUrl, "b", "http://localhost:8080", "-b <base address result>")
}

func main() {
	flag.Parse()

	servAddr, exist := os.LookupEnv("USERNAME")
	if exist {
		// Print the value of the environment variable
		fmt.Print(servAddr)
	}

	baseUrl, exist := os.LookupEnv("BASE_URL")
	if exist {
		// Print the value of the environment variable
		fmt.Print(baseUrl)
	}

	e := echo.New()

	decodeHandler := handler.New(&options.resUrl)
	e.POST(`/`, decodeHandler.EncodeUrl)
	e.GET(`/{id}`, decodeHandler.DecodeUrl)

	defaultPort := "5000"

	if &options.port != nil {
		defaultPort = options.port
	}

	e.Logger.Fatal(e.Start(":" + defaultPort))
}
