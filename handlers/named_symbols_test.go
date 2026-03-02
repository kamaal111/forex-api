package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamaal111/forex-api/utils"
)

func TestableNamedSymbolsHandler(repo RatesRepository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		service := NewRatesService(repo)

		record, err := service.GetAllNamedSymbols()
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if record == nil {
			utils.ErrorHandler(writer, "symbols not found", http.StatusNotFound)
			return
		}

		output, err := json.Marshal(record)
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("content-type", "application/json")
		writer.Write(output)
	}
}

func TestGetNamedSymbolsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockRecord     *SymbolsRecord
		mockErr        error
		wantStatusCode int
		wantSymbols    []NamedSymbol
	}{
		{
			name:           "returns named symbols for symbols that have rates in the database",
			mockRecord:     &SymbolsRecord{Date: "2025-11-21", Symbols: []string{"EUR", "USD", "GBP"}},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSymbols: []NamedSymbol{
				{Symbol: "EUR", Name: "Euro"},
				{Symbol: "USD", Name: "US Dollar"},
				{Symbol: "GBP", Name: "British Pound Sterling"},
			},
		},
		{
			name:           "returns 404 when no data exists in the database",
			mockRecord:     nil,
			mockErr:        nil,
			wantStatusCode: http.StatusNotFound,
			wantSymbols:    nil,
		},
		{
			name:           "returns 500 on database error",
			mockRecord:     nil,
			mockErr:        errors.New("database error"),
			wantStatusCode: http.StatusInternalServerError,
			wantSymbols:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetAllSymbolsFunc: func() (*SymbolsRecord, error) {
					return tt.mockRecord, tt.mockErr
				},
			}

			handler := TestableNamedSymbolsHandler(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/v1/rates/named-symbols", nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("GetNamedSymbols() status = %d, want %d", recorder.Code, tt.wantStatusCode)
			}

			if tt.wantSymbols != nil {
				contentType := recorder.Header().Get("content-type")
				if contentType != "application/json" {
					t.Errorf("GetNamedSymbols() content-type = %q, want %q", contentType, "application/json")
				}

				var record NamedSymbolsRecord
				if err := json.NewDecoder(recorder.Body).Decode(&record); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(record.Symbols) != len(tt.wantSymbols) {
					t.Errorf("GetNamedSymbols() returned %d symbols, want %d", len(record.Symbols), len(tt.wantSymbols))
				}

				for i, expected := range tt.wantSymbols {
					if i >= len(record.Symbols) {
						break
					}
					if record.Symbols[i].Symbol != expected.Symbol {
						t.Errorf("GetNamedSymbols() symbols[%d].symbol = %q, want %q", i, record.Symbols[i].Symbol, expected.Symbol)
					}
					if record.Symbols[i].Name != expected.Name {
						t.Errorf("GetNamedSymbols() symbols[%d].name = %q, want %q", i, record.Symbols[i].Name, expected.Name)
					}
				}
			}
		})
	}
}
