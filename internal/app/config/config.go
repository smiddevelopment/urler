package config

import (
	"flag"
	"os"
)

type netAddress struct {
	ServAddr string
	ResURL   string
}

var NetAddress netAddress

func SetConfig() {
	flag.String("a", "localhost:8080", "-a server address")
	flag.String("b", "http://localhost:8080", "-b result URL address")
	flag.Parse()

	NetAddress = netAddress{
		ServAddr: "localhost:8080",
		ResURL:   "http://localhost:8080",
	}

	if flag.Lookup("a") != nil {
		NetAddress.ServAddr = flag.Lookup("a").Value.(flag.Getter).String()
	}

	if flag.Lookup("b") != nil {
		NetAddress.ResURL = flag.Lookup("b").Value.(flag.Getter).String()
	}

	servURL, exist := os.LookupEnv("SERVER_ADDRESS")
	if exist {
		NetAddress.ServAddr = servURL
	}

	resURL, exist := os.LookupEnv("BASE_URL")
	if exist {
		NetAddress.ResURL = resURL
	}
}
