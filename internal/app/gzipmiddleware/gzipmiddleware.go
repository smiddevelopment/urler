package gzipmiddleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"golang.org/x/exp/slices"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// GzipMiddleware принимает параметром Handler и возвращает тоже Handler.
func GzipMiddleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(r.Header.Values("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return

		}

		if r.Header.Get("Content-Type") != "text/plain" && r.Header.Get("Content-Type") != "application/json" {
			next.ServeHTTP(w, r)
			return

		}

		if slices.Contains(r.Header.Values("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}

			body, _ := io.ReadAll(gz)

			defer gz.Close()

			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r)

		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return

		}
		_ = gz.Flush()
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
