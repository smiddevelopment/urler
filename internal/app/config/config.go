package config

import (
	"flag"
	"os"
)

type serverConfig struct {
	ServAddr string
	ResURL   string
	URLFile  string
	DBURL    string
}

var ServerConfig serverConfig

func SetConfig() {
	flag.String("a", "localhost:8080", "-a server address")
	flag.String("b", "http://localhost:8080", "-b result URL address")
	pwd, _ := os.Getwd()
	flag.String("f", pwd+"/cmd/shortener/short-url-db.json", "-f urls file db")
	flag.String("d", "user=smiddevelopment password=d5mvZQ2CTDYs dbname=gotest host=ep-round-sea-739368.eu-central-1.aws.neon.tech sslmode=verify-full", "-d url to db")
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

	if flag.Lookup("d") != nil {
		ServerConfig.DBURL = flag.Lookup("d").Value.(flag.Getter).String()
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

	dbURL, exist := os.LookupEnv("DATABASE_DSN")
	if exist {
		ServerConfig.DBURL = dbURL
	}
}
