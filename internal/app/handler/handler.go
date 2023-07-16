package handler

import (
	"fmt"
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
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "30")
		w.WriteHeader(http.StatusCreated)

		_, err := w.Write([]byte("http://localhost:8080/" + shortener.EncodeString(bodyString)))
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

	for _, v := range r.Cookies() {
		fmt.Printf("%s = %s\r\n", v.Name, v.Value)
	}

	if r.URL.Path == "" {
		http.Error(w, "id is empty!", http.StatusBadRequest)

		return
	}
	resLink := shortener.DecodeString(strings.TrimPrefix(r.URL.Path, "/"))
	if resLink != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", resLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

	http.Error(w, "this Id invalid!", http.StatusBadRequest)
}
