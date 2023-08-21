package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/smiddevelopment/urler.git/internal/app/gzipmiddleware"

	"github.com/smiddevelopment/urler.git/internal/app/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeURLHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "encode url #1",
			URL:  "https://practicum.yandex.ru/",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.URL))
			request.Header.Add("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			EncodeURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			resBody, err := io.ReadAll(res.Body)
			// Отложенное особождение памяти
			defer res.Body.Close()

			require.NoError(t, err)

			if string(resBody) == "" {
				t.Errorf("EncodeURL() = resBody is empty!")
			}

			if len(storage.EncodedURLs) == 0 {
				t.Errorf("EncodedURLs is empty!")
			}

			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestEncodeURLJSONHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "encode url #1",
			URL:  `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(test.URL)))
			request.Header.Add("Content-Type", "application/json")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			EncodeURLJSON(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			var getURL storage.URLEncoded
			if err := json.NewDecoder(res.Body).Decode(&getURL); err != nil {
				t.Errorf("NewDecoder() = " + err.Error())
				return

			}

			// Отложенное особождение памяти
			defer res.Body.Close()

			if string(getURL.Result) == "" {
				t.Errorf("EncodeURLJSON() = resBody is empty!")
			}

			if len(storage.EncodedURLs) == 0 {
				t.Errorf("EncodedURLs is empty!")
			}

			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestDecodeUrlHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "decode url #1",
			want: want{
				code:        307,
				response:    "Invalid ID!",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Add("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			DecodeURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			// Особождение памяти
			if err := res.Body.Close(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
			if res.Header.Get("Location") == "" {
				t.Errorf("There is no 'Location' header!")
			}
			assert.Equal(t, res.Header.Get("Location"), test.want.response)
		})
	}
}

func TestIntegralGZIPEncodeURLJSONHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "encode url #1",
			URL:  `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var b = []byte(test.URL)
			var buf bytes.Buffer
			gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
			if err != nil {
				_, _ = io.WriteString(gz, err.Error())
				return

			}
			if _, err = gz.Write(b); err != nil {
				log.Fatal(err)
				return
			}
			_ = gz.Flush()
			defer gz.Close()

			r := chi.NewRouter()
			r.Use(gzipmiddleware.GzipMiddleware)
			r.Post("/api/shorten", EncodeURLJSON)

			ts := httptest.NewServer(r)
			defer ts.Close()

			request, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", &buf)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept-Encoding", "gzip")
			request.Header.Set("Content-Encoding", "gzip")
			// создаём новый Recorder
			//w := httptest.NewRecorder()

			res, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			//res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			var getURL storage.URLEncoded
			if err := json.NewDecoder(res.Body).Decode(&getURL); err != nil {
				t.Errorf("NewDecoder() = " + err.Error())
				return

			}

			// Отложенное особождение памяти
			defer res.Body.Close()

			if string(getURL.Result) == "" {
				t.Errorf("EncodeURLJSON() = resBody is empty!")
			}

			if len(storage.EncodedURLs) == 0 {
				t.Errorf("EncodedURLs is empty!")
			}

			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
