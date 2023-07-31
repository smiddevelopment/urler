package config

import (
	"flag"
	"os"
)

type serverConfig struct {
	ServAddr string
	ResURL   string
	URLFile  string
}

var ServerConfig serverConfig

func SetConfig() {
	flag.String("a", "localhost:8080", "-a server address")
	flag.String("b", "http://localhost:8080", "-b result URL address")
	pwd, _ := os.Getwd()
	flag.String("f", pwd+"/cmd/shortener/short-url-db.json", "-b urls file db")
	flag.Parse()

	ServerConfig = serverConfig{
		ServAddr: "localhost:8080",
		ResURL:   "http://localhost:8080",
	}

	if flag.Lookup("a") != nil {
		ServerConfig.ServAddr = flag.Lookup("a").Value.(flag.Getter).String()
	}

	if flag.Lookup("b") != nil {
		ServerConfig.ResURL = flag.Lookup("b").Value.(flag.Getter).String()
	}

	if flag.Lookup("f") != nil {
		ServerConfig.URLFile = flag.Lookup("f").Value.(flag.Getter).String()
	}

	servURL, exist := os.LookupEnv("SERVER_ADDRESS")
	if exist {
		ServerConfig.ServAddr = servURL
	}

	resURL, exist := os.LookupEnv("BASE_URL")
	if exist {
		ServerConfig.ResURL = resURL
	}

	uRLFile, exist := os.LookupEnv("FILE_STORAGE_PATH")
	if exist {
		ServerConfig.URLFile = uRLFile
	}
}
