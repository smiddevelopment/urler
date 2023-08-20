package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/smiddevelopment/urler.git/internal/app/config"

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
	defer r.Body.Close()

	bodyString := string(body)
	if bodyString != "" {
		resURL := config.ServerConfig.ResURL + "/" + storage.Add(bodyString)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len([]rune(resURL))))
		w.WriteHeader(http.StatusCreated)
		// Получение значения ID из хранилища или добавление новой ссылки
		_, err := w.Write([]byte(resURL))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	http.Error(w, "body is empty!", http.StatusBadRequest)
}

// EncodeURLJSON обработка запроса POST, кодирование ссылки
func EncodeURLJSON(w http.ResponseWriter, r *http.Request) {
	// Чтение тела запроса
	var getURL storage.URLEncoded
	var sendURL storage.URLEncoded

	if err := json.NewDecoder(r.Body).Decode(&getURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	// Отложенное особождение памяти
	defer r.Body.Close()

	sendURL.Result = config.ServerConfig.ResURL + "/" + storage.Add(getURL.URL)
	w.Header().Set("Content-Type", "application/json")
	stringJSON, marshalErr := json.Marshal(sendURL)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusBadRequest)
		return

	}

	w.Header().Set("Content-Length", strconv.Itoa(len(string(stringJSON))+1))
	w.WriteHeader(http.StatusCreated)

	// Получение значения ID из хранилища или добавление новой ссылки
	err := json.NewEncoder(w).Encode(sendURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
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

	http.Error(w, "this Id invalid!", http.StatusNotFound)
}

// PingDB проверка подключения к базе данных
func PingDB(w http.ResponseWriter, r *http.Request) {
	// Попытка установить соединение с базой данных
	db, err := sql.Open("pgx", config.ServerConfig.DbURL)
	if err != nil {
		http.Error(w, "Can't connect to database!", http.StatusInternalServerError)
	}
	defer db.Close()
	w.WriteHeader(http.StatusOK)

	return
}
