package middleware

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Ключ для безопасного доступа к контексту
type contextKey string

const requestIDKey contextKey = "request_id"

// Генерация случайного ID запроса
func generateRequestID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RequestIDMiddleware Middleware для обработки request_id
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем или генерируем request_id
		requestID := r.URL.Query().Get("request_id")
		if requestID == "" {
			requestID = generateRequestID(12)
		}

		// Добавляем в контекст с безопасным ключом
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Кастомный ResponseWriter для перехвата статуса ответа
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware Middleware для логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Используем кастомный ResponseWriter
		rw := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

		// Выполняем цепочку обработчиков
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Извлекаем request_id из контекста
		requestID := ""
		if id, ok := r.Context().Value(requestIDKey).(string); ok {
			requestID = id
		}

		// Формируем лог
		log.Printf(
			"[%s] %s %s %d %v %s",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			rw.status,
			duration,
			requestID,
		)
	})
}
