package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		code        int
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "not found error",
			message:     "Resource not found",
			code:        http.StatusNotFound,
			wantStatus:  http.StatusNotFound,
			wantMessage: "Resource not found",
		},
		{
			name:        "internal server error",
			message:     "Something went wrong",
			code:        http.StatusInternalServerError,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Something went wrong",
		},
		{
			name:        "bad request error",
			message:     "Invalid input",
			code:        http.StatusBadRequest,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Invalid input",
		},
		{
			name:        "unauthorized error",
			message:     "Unauthorized access",
			code:        http.StatusUnauthorized,
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Unauthorized access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			ErrorHandler(recorder, tt.message, tt.code)

			if recorder.Code != tt.wantStatus {
				t.Errorf("ErrorHandler() status = %d, want %d", recorder.Code, tt.wantStatus)
			}

			contentType := recorder.Header().Get("content-type")
			if contentType != "application/json" {
				t.Errorf("ErrorHandler() content-type = %q, want %q", contentType, "application/json")
			}

			var gotError Error
			if err := json.NewDecoder(recorder.Body).Decode(&gotError); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if gotError.Message != tt.wantMessage {
				t.Errorf("ErrorHandler() message = %q, want %q", gotError.Message, tt.wantMessage)
			}

			if gotError.Status != tt.wantStatus {
				t.Errorf("ErrorHandler() status in body = %d, want %d", gotError.Status, tt.wantStatus)
			}
		})
	}
}
