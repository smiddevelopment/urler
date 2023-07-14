package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeUrlHandler(t *testing.T) {
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
			name: "encode url #1",
			want: want{
				code:        201,
				response:    "EwHXdJfB",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/"))
			request.Header.Add("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			EncodeUrl(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, string(resBody), test.want.response)
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
			name: "decode url test #1",
			want: want{
				code:        307,
				response:    "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodGet, "/EwHXdJfB", nil)
			request.Header.Add("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			DecodeUrl(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
			assert.Equal(t, res.Header.Get("Location"), test.want.response)
		})
	}
}
