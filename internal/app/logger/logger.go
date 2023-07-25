package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func InitLog() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()
}

type DataWriter struct {
	http.ResponseWriter

	status int
	size   int
}

func (w *DataWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *DataWriter) Write(b []byte) (int, error) {
	w.size = len(b)
	return w.ResponseWriter.Write(b)
}

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		dw := &DataWriter{ResponseWriter: w}
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// точка, где выполняется хендлер
		h.ServeHTTP(dw, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

		// отправляем сведения об ответе в zap
		sugar.Infoln(
			"status", dw.status,
			"size", dw.size,
		)

	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
