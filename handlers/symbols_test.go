package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamaal111/forex-api/utils"
)

func TestableSymbolsHandler(repo RatesRepository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		service := NewRatesService(repo)

		symbols, err := service.GetAllSymbols()
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		output, err := json.Marshal(symbols)
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("content-type", "application/json")
		writer.Write(output)
	}
}

func TestGetSymbolsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockSymbols    []string
		mockErr        error
		wantStatusCode int
		wantSymbols    []string
	}{
		{
			name:           "returns only symbols that have rates in the database",
			mockSymbols:    []string{"EUR", "USD", "GBP"},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSymbols:    []string{"EUR", "USD", "GBP"},
		},
		{
			name:           "returns empty list when no data exists in the database",
			mockSymbols:    []string{},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSymbols:    []string{},
		},
		{
			name:           "returns 500 on database error",
			mockSymbols:    nil,
			mockErr:        errors.New("database error"),
			wantStatusCode: http.StatusInternalServerError,
			wantSymbols:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetAllSymbolsFunc: func() ([]string, error) {
					return tt.mockSymbols, tt.mockErr
				},
			}

			handler := TestableSymbolsHandler(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/v1/rates/symbols", nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("GetSymbols() status = %d, want %d", recorder.Code, tt.wantStatusCode)
			}

			contentType := recorder.Header().Get("content-type")
			if contentType != "application/json" {
				t.Errorf("GetSymbols() content-type = %q, want %q", contentType, "application/json")
			}

			if tt.wantStatusCode == http.StatusOK {
				var symbols []string
				if err := json.NewDecoder(recorder.Body).Decode(&symbols); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(symbols) != len(tt.wantSymbols) {
					t.Errorf("GetSymbols() returned %d symbols, want %d", len(symbols), len(tt.wantSymbols))
				}

				symbolSet := make(map[string]bool, len(symbols))
				for _, s := range symbols {
					symbolSet[s] = true
				}
				for _, expected := range tt.wantSymbols {
					if !symbolSet[expected] {
						t.Errorf("GetSymbols() missing expected symbol %q", expected)
					}
				}
			}
		})
	}
}
