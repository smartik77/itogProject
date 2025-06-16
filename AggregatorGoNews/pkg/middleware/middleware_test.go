package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	t.Run("generates ID of correct length", func(t *testing.T) {
		id := generateRequestID(12)
		if len(id) != 12 {
			t.Errorf("Expected length 12, got %d", len(id))
		}
	})

	t.Run("generates different IDs", func(t *testing.T) {
		id1 := generateRequestID(12)
		id2 := generateRequestID(12)
		if id1 == id2 {
			t.Error("Expected different IDs, got the same")
		}
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("uses existing request_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?request_id=test123", nil)
		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(requestIDKey).(string)
			if id != "test123" {
				t.Errorf("Expected request_id 'test123', got '%s'", id)
			}
		})

		middleware := RequestIDMiddleware(testHandler)
		middleware.ServeHTTP(rr, req)
	})

	t.Run("generates new request_id if missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(requestIDKey).(string)
			if id == "" {
				t.Error("Expected generated request_id, got empty")
			}
		})

		middleware := RequestIDMiddleware(testHandler)
		middleware.ServeHTTP(rr, req)
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("logs request information", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		// Перехватываем вывод логов
		var logOutput strings.Builder
		log.SetOutput(&logOutput)
		defer log.SetOutput(nil)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := LoggingMiddleware(testHandler)
		middleware.ServeHTTP(rr, req)

		logStr := logOutput.String()
		if !strings.Contains(logStr, "GET /test 200") {
			t.Errorf("Log output missing expected data: %s", logStr)
		}
	})
}

func TestResponseRecorder(t *testing.T) {
	t.Run("records status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rec := &responseRecorder{ResponseWriter: rr}

		rec.WriteHeader(http.StatusNotFound)
		if rec.status != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.status)
		}
	})
}
