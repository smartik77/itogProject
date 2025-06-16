package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	t.Run("should generate ID of correct length", func(t *testing.T) {
		id := generateRequestID(12)
		if len(id) != 12 {
			t.Errorf("Expected ID length 12, got %d", len(id))
		}
	})

	t.Run("should generate different IDs each time", func(t *testing.T) {
		id1 := generateRequestID(12)
		id2 := generateRequestID(12)
		if id1 == id2 {
			t.Error("Expected different IDs, got the same")
		}
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should use existing request ID from query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?request_id=existing123", nil)
		rr := httptest.NewRecorder()

		var capturedID string
		handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedID = r.Context().Value(requestIDKey).(string)
		}))

		handler.ServeHTTP(rr, req)

		if capturedID != "existing123" {
			t.Errorf("Expected request ID 'existing123', got '%s'", capturedID)
		}
	})

	t.Run("should generate new request ID if not provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		var capturedID string
		handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedID = r.Context().Value(requestIDKey).(string)
		}))

		handler.ServeHTTP(rr, req)

		if capturedID == "" || len(capturedID) != 12 {
			t.Errorf("Expected generated 12-char ID, got '%s'", capturedID)
		}
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("should log request details", func(t *testing.T) {
		var logOutput bytes.Buffer
		log.SetOutput(&logOutput)
		defer log.SetOutput(nil)

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		rr := httptest.NewRecorder()

		// Добавляем request ID в контекст
		ctx := context.WithValue(req.Context(), requestIDKey, "test123")
		req = req.WithContext(ctx)

		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)

		logStr := logOutput.String()
		expectedParts := []string{
			"GET /test",
			"200",
			"test123",
			"127.0.0.1:12345",
		}

		for _, part := range expectedParts {
			if !strings.Contains(logStr, part) {
				t.Errorf("Expected log to contain '%s', got '%s'", part, logStr)
			}
		}
	})

	t.Run("should record response status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
		}
	})
}

func TestResponseRecorder(t *testing.T) {
	t.Run("should record status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rec := &responseRecorder{ResponseWriter: rr, status: http.StatusOK}

		testStatus := http.StatusNotFound
		rec.WriteHeader(testStatus)

		if rec.status != testStatus {
			t.Errorf("Expected status %d, got %d", testStatus, rec.status)
		}
		if rr.Code != testStatus {
			t.Errorf("Expected underlying writer status %d, got %d", testStatus, rr.Code)
		}
	})
}
