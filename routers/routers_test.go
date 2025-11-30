package routers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamaal111/forex-api/utils"
)

func TestNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	recorder := httptest.NewRecorder()

	notFound(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("notFound() status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	contentType := recorder.Header().Get("content-type")
	if contentType != "application/json" {
		t.Errorf("notFound() content-type = %q, want %q", contentType, "application/json")
	}

	var response utils.Error
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Message != "Not found" {
		t.Errorf("notFound() message = %q, want %q", response.Message, "Not found")
	}

	if response.Status != http.StatusNotFound {
		t.Errorf("notFound() status in body = %d, want %d", response.Status, http.StatusNotFound)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := loggerMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("loggerMiddleware() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Body.String() != "OK" {
		t.Errorf("loggerMiddleware() body = %q, want %q", recorder.Body.String(), "OK")
	}
}

func TestLoggerMiddleware_CapturesStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		wantStatusCode int
	}{
		{
			name:           "captures 200 OK",
			statusCode:     http.StatusOK,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "captures 404 Not Found",
			statusCode:     http.StatusNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "captures 500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "captures 201 Created",
			statusCode:     http.StatusCreated,
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			handler := loggerMiddleware(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("loggerMiddleware() status = %d, want %d", recorder.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestResponseObserver_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	observer := &responseObserver{ResponseWriter: recorder}

	observer.WriteHeader(http.StatusCreated)

	if observer.status != http.StatusCreated {
		t.Errorf("responseObserver.status = %d, want %d", observer.status, http.StatusCreated)
	}

	if recorder.Code != http.StatusCreated {
		t.Errorf("underlying ResponseWriter status = %d, want %d", recorder.Code, http.StatusCreated)
	}
}

func TestResponseObserver_Write(t *testing.T) {
	recorder := httptest.NewRecorder()
	observer := &responseObserver{ResponseWriter: recorder}

	content := []byte("test content")
	n, err := observer.Write(content)

	if err != nil {
		t.Errorf("responseObserver.Write() error = %v", err)
	}

	if n != len(content) {
		t.Errorf("responseObserver.Write() returned %d bytes, want %d", n, len(content))
	}

	if recorder.Body.String() != "test content" {
		t.Errorf("responseObserver.Write() body = %q, want %q", recorder.Body.String(), "test content")
	}
}
