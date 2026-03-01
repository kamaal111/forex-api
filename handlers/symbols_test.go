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

		record, err := service.GetAllSymbols()
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

func TestGetSymbolsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockRecord     *SymbolsRecord
		mockErr        error
		wantStatusCode int
		wantRecord     *SymbolsRecord
	}{
		{
			name:           "returns only symbols that have rates in the database",
			mockRecord:     &SymbolsRecord{Date: "2025-11-21", Symbols: []string{"EUR", "USD", "GBP"}},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantRecord:     &SymbolsRecord{Date: "2025-11-21", Symbols: []string{"EUR", "USD", "GBP"}},
		},
		{
			name:           "returns 404 when no data exists in the database",
			mockRecord:     nil,
			mockErr:        nil,
			wantStatusCode: http.StatusNotFound,
			wantRecord:     nil,
		},
		{
			name:           "returns 500 on database error",
			mockRecord:     nil,
			mockErr:        errors.New("database error"),
			wantStatusCode: http.StatusInternalServerError,
			wantRecord:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetAllSymbolsFunc: func() (*SymbolsRecord, error) {
					return tt.mockRecord, tt.mockErr
				},
			}

			handler := TestableSymbolsHandler(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/v1/rates/symbols", nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("GetSymbols() status = %d, want %d", recorder.Code, tt.wantStatusCode)
			}

			if tt.wantRecord != nil {
				contentType := recorder.Header().Get("content-type")
				if contentType != "application/json" {
					t.Errorf("GetSymbols() content-type = %q, want %q", contentType, "application/json")
				}

				var record SymbolsRecord
				if err := json.NewDecoder(recorder.Body).Decode(&record); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if record.Date != tt.wantRecord.Date {
					t.Errorf("GetSymbols() date = %q, want %q", record.Date, tt.wantRecord.Date)
				}

				if len(record.Symbols) != len(tt.wantRecord.Symbols) {
					t.Errorf("GetSymbols() returned %d symbols, want %d", len(record.Symbols), len(tt.wantRecord.Symbols))
				}

				for i, expected := range tt.wantRecord.Symbols {
					if i >= len(record.Symbols) {
						break
					}
					if record.Symbols[i] != expected {
						t.Errorf("GetSymbols() symbols[%d] = %q, want %q", i, record.Symbols[i], expected)
					}
				}
			}
		})
	}
}
