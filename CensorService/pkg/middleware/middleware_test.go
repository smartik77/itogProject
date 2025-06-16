package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	t.Run("generates ID of specified length", func(t *testing.T) {
		id := generateRequestID(15)
		if len(id) != 15 {
			t.Errorf("Expected length 15, got %d", len(id))
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := generateRequestID(10)
		id2 := generateRequestID(10)
		if id1 == id2 {
			t.Error("Expected different IDs, got the same")
		}
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("uses existing request ID from query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/path?request_id=test123", nil)
		rr := httptest.NewRecorder()

		handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверка установки заголовка
			if id := w.Header().Get("X-Request-ID"); id != "test123" {
				t.Errorf("Header X-Request-ID = %v, want test123", id)
			}
		}))

		handler.ServeHTTP(rr, req)
	})

	t.Run("generates new request ID when missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/path", nil)
		rr := httptest.NewRecorder()

		handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := w.Header().Get("X-Request-ID")
			if len(id) != 12 {
				t.Errorf("Generated ID length = %d, want 12", len(id))
			}
		}))

		handler.ServeHTTP(rr, req)
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("logs request details correctly", func(t *testing.T) {
		// Перехватываем вывод логов
		var logOutput bytes.Buffer
		log.SetOutput(&logOutput)
		defer log.SetOutput(nil)

		req := httptest.NewRequest("POST", "/api", nil)
		rr := httptest.NewRecorder()

		// Устанавливаем request ID
		req.Header.Set("X-Request-ID", "log-test-123")

		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))

		handler.ServeHTTP(rr, req)

		logStr := logOutput.String()
		expected := []string{
			"POST /api",
			"201",
			"log-test-123",
		}

		for _, s := range expected {
			if !strings.Contains(logStr, s) {
				t.Errorf("Expected log to contain %q, got %q", s, logStr)
			}
		}
	})
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Run("records status code correctly", func(t *testing.T) {
		rr := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: rr}

		lrw.WriteHeader(http.StatusNotFound)
		if lrw.statusCode != http.StatusNotFound {
			t.Errorf("statusCode = %d, want %d", lrw.statusCode, http.StatusNotFound)
		}
		if rr.Code != http.StatusNotFound {
			t.Errorf("Underlying writer code = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("defaults to status OK", func(t *testing.T) {
		lrw := &loggingResponseWriter{ResponseWriter: httptest.NewRecorder()}
		if lrw.statusCode != 0 { // Проверяем начальное значение
			t.Errorf("Initial statusCode = %d, want 0", lrw.statusCode)
		}

		// При первом вызове WriteHeader устанавливается статус
		lrw.WriteHeader(http.StatusOK)
		if lrw.statusCode != http.StatusOK {
			t.Errorf("statusCode = %d, want %d", lrw.statusCode, http.StatusOK)
		}
	})
}
