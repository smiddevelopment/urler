package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/smiddevelopment/urler.git/internal/app/shortener"
)

func EncodeURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}(r.Body)
	bodyString := string(body)
	if bodyString != "" {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "30")

		_, err := w.Write([]byte(r.Host + "/" + shortener.EncodeString(bodyString)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	http.Error(w, "body is empty!", http.StatusBadRequest)
}

func DecodeURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if r.URL.Path == "" {
		http.Error(w, "id is empty!", http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", shortener.DecodeString(strings.TrimPrefix(r.URL.Path, "/")))
}
