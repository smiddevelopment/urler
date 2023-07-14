package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/smiddevelopment/urler.git/internal/app/shortener"
)

func EncodeUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()
	bodyString := string(body)
	if bodyString != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "30")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortener.EncodeString(bodyString)))

		return
	}

	http.Error(w, "body is empty!", http.StatusBadRequest)
}

func DecodeUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if r.URL.Path == "" {
		http.Error(w, "id is empty!", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", shortener.DecodeString(strings.TrimPrefix(r.URL.Path, "/")))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
