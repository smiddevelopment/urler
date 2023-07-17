package handler

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/smiddevelopment/urler.git/internal/app/storage"
)

// EncodeURL обработка запроса POST, кодирование ссылки
func EncodeURL(w http.ResponseWriter, r *http.Request) {
	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	// Отложенное особождение памяти
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
		// Получение значения ID из хранилища или добавление новой ссылки
		resURL, exist := os.LookupEnv("BASE_URL")
		if exist {
			fmt.Print(resURL)
		}

		if resURL == "" {
			resURL = flag.Lookup("b").Value.(flag.Getter).String()
		}

		if resURL == "" {
			resURL = "http://localhost:8080"
		}

		_, err := w.Write([]byte(resURL + "/" + storage.Add(bodyString)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	http.Error(w, "body is empty!", http.StatusBadRequest)
}

// DecodeURL обработка запроса GET, декодирование ссылки
func DecodeURL(w http.ResponseWriter, r *http.Request) {
	// Получение значения URL из хранилища, если найдено
	resLink := storage.Get(strings.TrimPrefix(r.URL.Path, "/"))
	if resLink != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", resLink)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	http.Error(w, "this Id invalid!", http.StatusBadRequest)
}
