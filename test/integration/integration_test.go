package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type ExchangeRateRecord struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func TestGetLatestEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := NewTestContext()
	if err := tc.Setup(0); err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tc.Teardown()

	t.Run("returns 404 when no data exists", func(t *testing.T) {
		if err := tc.ClearCollection("exchange_rates"); err != nil {
			t.Fatalf("Failed to clear collection: %v", err)
		}

		resp, err := tc.Server.GetLatest("", "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns data with default base EUR", func(t *testing.T) {
		if err := tc.ClearCollection("exchange_rates"); err != nil {
			t.Fatalf("Failed to clear collection: %v", err)
		}

		_, err := tc.DB.Collection("exchange_rates").Doc("EUR-2025-11-21").Set(tc.Ctx, map[string]interface{}{
			"base": "EUR",
			"date": "2025-11-21",
			"rates": map[string]float64{
				"USD": 1.08,
				"GBP": 0.86,
				"JPY": 161.5,
			},
		})
		if err != nil {
			t.Fatalf("Failed to seed data: %v", err)
		}

		resp, err := tc.Server.GetLatest("", "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		var record ExchangeRateRecord
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if record.Base != "EUR" {
			t.Errorf("Expected base EUR, got %s", record.Base)
		}
		if record.Date != "2025-11-21" {
			t.Errorf("Expected date 2025-11-21, got %s", record.Date)
		}
		if len(record.Rates) != 3 {
			t.Errorf("Expected 3 rates, got %d", len(record.Rates))
		}
	})

	t.Run("returns data with specific base", func(t *testing.T) {
		if err := tc.ClearCollection("exchange_rates"); err != nil {
			t.Fatalf("Failed to clear collection: %v", err)
		}

		_, err := tc.DB.Collection("exchange_rates").Doc("USD-2025-11-21").Set(tc.Ctx, map[string]interface{}{
			"base": "USD",
			"date": "2025-11-21",
			"rates": map[string]float64{
				"EUR": 0.926,
				"GBP": 0.796,
				"JPY": 149.5,
			},
		})
		if err != nil {
			t.Fatalf("Failed to seed data: %v", err)
		}

		resp, err := tc.Server.GetLatest("USD", "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		var record ExchangeRateRecord
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if record.Base != "USD" {
			t.Errorf("Expected base USD, got %s", record.Base)
		}
	})

	t.Run("filters rates by symbols", func(t *testing.T) {
		if err := tc.ClearCollection("exchange_rates"); err != nil {
			t.Fatalf("Failed to clear collection: %v", err)
		}

		_, err := tc.DB.Collection("exchange_rates").Doc("EUR-2025-11-21").Set(tc.Ctx, map[string]any{
			"base": "EUR",
			"date": "2025-11-21",
			"rates": map[string]float64{
				"USD": 1.08,
				"GBP": 0.86,
				"JPY": 161.5,
				"CHF": 0.93,
			},
		})
		if err != nil {
			t.Fatalf("Failed to seed data: %v", err)
		}

		resp, err := tc.Server.GetLatest("EUR", "USD,GBP")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		var record ExchangeRateRecord
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(record.Rates) != 2 {
			t.Errorf("Expected 2 rates, got %d", len(record.Rates))
		}

		if _, ok := record.Rates["USD"]; !ok {
			t.Error("Expected USD in rates")
		}
		if _, ok := record.Rates["GBP"]; !ok {
			t.Error("Expected GBP in rates")
		}
		if _, ok := record.Rates["JPY"]; ok {
			t.Error("Did not expect JPY in filtered rates")
		}
	})

	t.Run("returns latest rate when multiple dates exist", func(t *testing.T) {
		if err := tc.ClearCollection("exchange_rates"); err != nil {
			t.Fatalf("Failed to clear collection: %v", err)
		}

		_, err := tc.DB.Collection("exchange_rates").Doc("EUR-2025-11-20").Set(tc.Ctx, map[string]interface{}{
			"base": "EUR",
			"date": "2025-11-20",
			"rates": map[string]float64{
				"USD": 1.07,
			},
		})
		if err != nil {
			t.Fatalf("Failed to seed old data: %v", err)
		}

		_, err = tc.DB.Collection("exchange_rates").Doc("EUR-2025-11-21").Set(tc.Ctx, map[string]any{
			"base": "EUR",
			"date": "2025-11-21",
			"rates": map[string]float64{
				"USD": 1.08,
			},
		})
		if err != nil {
			t.Fatalf("Failed to seed new data: %v", err)
		}

		resp, err := tc.Server.GetLatest("EUR", "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		var record ExchangeRateRecord
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if record.Date != "2025-11-21" {
			t.Errorf("Expected latest date 2025-11-21, got %s", record.Date)
		}
		if record.Rates["USD"] != 1.08 {
			t.Errorf("Expected USD rate 1.08, got %f", record.Rates["USD"])
		}
	})
}

func TestContentType(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := NewTestContext()
	if err := tc.Setup(0); err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tc.Teardown()

	_, err := tc.DB.Collection("exchange_rates").Doc("EUR-2025-11-21").Set(tc.Ctx, map[string]interface{}{
		"base": "EUR",
		"date": "2025-11-21",
		"rates": map[string]float64{
			"USD": 1.08,
		},
	})
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	resp, err := tc.Server.GetLatest("", "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}
