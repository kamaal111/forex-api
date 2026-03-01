package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSymbolsHandler(t *testing.T) {
	t.Run("returns 200 with all supported currency symbols", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/rates/symbols", nil)
		recorder := httptest.NewRecorder()

		GetSymbols(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("GetSymbols() status = %d, want %d", recorder.Code, http.StatusOK)
		}

		contentType := recorder.Header().Get("content-type")
		if contentType != "application/json" {
			t.Errorf("GetSymbols() content-type = %q, want %q", contentType, "application/json")
		}

		var symbols []string
		if err := json.NewDecoder(recorder.Body).Decode(&symbols); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(symbols) != len(Currencies) {
			t.Errorf("GetSymbols() returned %d symbols, want %d", len(symbols), len(Currencies))
		}

		symbolSet := make(map[string]bool, len(symbols))
		for _, s := range symbols {
			symbolSet[s] = true
		}
		for _, expected := range Currencies {
			if !symbolSet[expected] {
				t.Errorf("GetSymbols() missing currency %q", expected)
			}
		}
	})

	t.Run("response contains well-known currencies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/rates/symbols", nil)
		recorder := httptest.NewRecorder()

		GetSymbols(recorder, req)

		var symbols []string
		if err := json.NewDecoder(recorder.Body).Decode(&symbols); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		wellKnown := []string{"EUR", "USD", "GBP", "JPY", "CHF", "AUD", "CAD"}
		symbolSet := make(map[string]bool, len(symbols))
		for _, s := range symbols {
			symbolSet[s] = true
		}
		for _, currency := range wellKnown {
			if !symbolSet[currency] {
				t.Errorf("GetSymbols() missing well-known currency %q", currency)
			}
		}
	})
}
