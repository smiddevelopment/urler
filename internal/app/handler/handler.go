package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/smiddevelopment/urler.git/internal/app/storage"
)

func RouteURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		DecodeURL(w, r)
	case http.MethodPost:
		EncodeURL(w, r)
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func EncodeURL(w http.ResponseWriter, r *http.Request) {
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

		_, err := w.Write([]byte("http://localhost:8080/" + storage.Add(bodyString)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	http.Error(w, "body is empty!", http.StatusBadRequest)
}

func DecodeURL(w http.ResponseWriter, r *http.Request) {
	resLink := storage.Get(strings.TrimPrefix(r.URL.Path, "/"))
	if resLink != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", resLink)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	http.Error(w, "this Id invalid!", http.StatusBadRequest)
}
