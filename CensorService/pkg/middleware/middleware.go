package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

// RequestIDMiddleware создает middleware для сквозного идентификатора запроса
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.URL.Query().Get("request_id")
		if requestID == "" {
			requestID = generateRequestID(12)
		}

		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware создает middleware для логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		requestID := w.Header().Get("X-Request-ID")

		log.Printf(
			"[%s] %s %s %d %v %s",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration,
			requestID,
		)
	})
}

// loggingResponseWriter кастомный ResponseWriter для перехвата статуса
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// generateRequestID генерирует случайный ID запроса
func generateRequestID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
