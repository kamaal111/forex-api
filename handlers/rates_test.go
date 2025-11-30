package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamaal111/forex-api/utils"
)

func TestableHandler(repo RatesRepository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		service := NewRatesService(repo)

		base := request.URL.Query().Get("base")
		symbols := request.URL.Query().Get("symbols")

		record, err := service.GetLatestRate(base, symbols)
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if record == nil {
			utils.ErrorHandler(writer, "Rates not found", http.StatusNotFound)
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

func TestGetLatestHandler(t *testing.T) {
	sampleRecord := &ExchangeRateRecord{
		Base: "EUR",
		Date: "2024-01-15",
		Rates: map[string]float64{
			"USD": 1.08,
			"GBP": 0.86,
			"JPY": 161.5,
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		mockRecord     *ExchangeRateRecord
		mockErr        error
		wantStatusCode int
		wantBase       string
		wantRatesCount int
	}{
		{
			name:           "successful request with default base",
			queryParams:    "",
			mockRecord:     sampleRecord,
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantBase:       "EUR",
			wantRatesCount: 3,
		},
		{
			name:           "successful request with explicit base",
			queryParams:    "?base=EUR",
			mockRecord:     sampleRecord,
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantBase:       "EUR",
			wantRatesCount: 3,
		},
		{
			name:           "successful request with symbols filter",
			queryParams:    "?base=EUR&symbols=USD",
			mockRecord:     sampleRecord,
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantBase:       "EUR",
			wantRatesCount: 1,
		},
		{
			name:           "not found when repository returns nil",
			queryParams:    "",
			mockRecord:     nil,
			mockErr:        nil,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "internal server error on repository error",
			queryParams:    "",
			mockRecord:     nil,
			mockErr:        ErrRatesNotFound,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetLatestRateFunc: func(base string) (*ExchangeRateRecord, error) {
					return tt.mockRecord, tt.mockErr
				},
			}

			handler := TestableHandler(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/v1/rates/latest"+tt.queryParams, nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("handler returned status %d, want %d", recorder.Code, tt.wantStatusCode)
			}

			contentType := recorder.Header().Get("content-type")
			if contentType != "application/json" {
				t.Errorf("handler returned content-type %q, want %q", contentType, "application/json")
			}

			if tt.wantStatusCode == http.StatusOK {
				var response ExchangeRateRecord
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response.Base != tt.wantBase {
					t.Errorf("response base = %q, want %q", response.Base, tt.wantBase)
				}

				if len(response.Rates) != tt.wantRatesCount {
					t.Errorf("response rates count = %d, want %d", len(response.Rates), tt.wantRatesCount)
				}
			}
		})
	}
}

func TestGetLatestHandler_SymbolsFiltering(t *testing.T) {
	sampleRecord := &ExchangeRateRecord{
		Base: "EUR",
		Date: "2024-01-15",
		Rates: map[string]float64{
			"USD": 1.08,
			"GBP": 0.86,
			"JPY": 161.5,
			"CHF": 0.94,
		},
	}

	tests := []struct {
		name        string
		symbols     string
		wantSymbols []string
	}{
		{
			name:        "single symbol",
			symbols:     "USD",
			wantSymbols: []string{"USD"},
		},
		{
			name:        "multiple symbols",
			symbols:     "USD,GBP",
			wantSymbols: []string{"USD", "GBP"},
		},
		{
			name:        "case insensitive",
			symbols:     "usd,gbp",
			wantSymbols: []string{"USD", "GBP"},
		},
		{
			name:        "excludes invalid symbols",
			symbols:     "USD,INVALID,GBP",
			wantSymbols: []string{"USD", "GBP"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetLatestRateFunc: func(base string) (*ExchangeRateRecord, error) {
					return sampleRecord, nil
				},
			}

			handler := TestableHandler(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/v1/rates/latest?symbols="+tt.symbols, nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			var response ExchangeRateRecord
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(response.Rates) != len(tt.wantSymbols) {
				t.Errorf("response rates count = %d, want %d", len(response.Rates), len(tt.wantSymbols))
			}

			for _, symbol := range tt.wantSymbols {
				if _, ok := response.Rates[symbol]; !ok {
					t.Errorf("response missing expected symbol: %s", symbol)
				}
			}
		})
	}
}
