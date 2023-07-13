package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/smiddevelopment/urler.git/internal/app/shortener"
)

func EncodeUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	bodyString := string(body)
	if bodyString != "" {
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(shortener.EncodeString(bodyString)))
	} else {
		http.Error(res, "Body is empty!", 500)
	}
}

func DecodeUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/")
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", shortener.DecodeString(id))
	res.WriteHeader(http.StatusTemporaryRedirect)
}
