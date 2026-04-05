package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetOpenAPISpecHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, OpenAPISpecPath, nil)
	recorder := httptest.NewRecorder()

	GetOpenAPISpec(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetOpenAPISpec() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/x-yaml" {
		t.Errorf("GetOpenAPISpec() Content-Type = %q, want %q", contentType, "application/x-yaml")
	}

	body := recorder.Body.String()
	if len(body) == 0 {
		t.Error("GetOpenAPISpec() returned empty body")
	}

	if !strings.Contains(body, "swagger:") && !strings.Contains(body, "openapi:") {
		t.Error("GetOpenAPISpec() response does not appear to be a valid OpenAPI spec")
	}
}
