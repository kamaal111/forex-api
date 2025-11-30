package handlers

import (
	"errors"
	"testing"
)

type MockRatesRepository struct {
	GetLatestRateFunc func(base string) (*ExchangeRateRecord, error)
}

func (m *MockRatesRepository) GetLatestRate(base string) (*ExchangeRateRecord, error) {
	if m.GetLatestRateFunc != nil {
		return m.GetLatestRateFunc(base)
	}
	return nil, nil
}

func TestNormalizeBase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid uppercase currency",
			input: "USD",
			want:  "USD",
		},
		{
			name:  "valid lowercase currency",
			input: "usd",
			want:  "USD",
		},
		{
			name:  "valid currency with whitespace",
			input: "  EUR  ",
			want:  "EUR",
		},
		{
			name:  "invalid currency defaults to EUR",
			input: "INVALID",
			want:  "EUR",
		},
		{
			name:  "empty string defaults to EUR",
			input: "",
			want:  "EUR",
		},
		{
			name:  "mixed case currency",
			input: "gBp",
			want:  "GBP",
		},
		{
			name:  "JPY currency",
			input: "jpy",
			want:  "JPY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeBase(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeBase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMakeSymbolsArray(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		base    string
		want    []string
		wantLen int
	}{
		{
			name:    "empty symbols",
			raw:     "",
			base:    "EUR",
			want:    []string{},
			wantLen: 0,
		},
		{
			name:    "single valid symbol",
			raw:     "USD",
			base:    "EUR",
			want:    []string{"USD"},
			wantLen: 1,
		},
		{
			name:    "multiple valid symbols",
			raw:     "USD,GBP,JPY",
			base:    "EUR",
			want:    []string{"USD", "GBP", "JPY"},
			wantLen: 3,
		},
		{
			name:    "symbols with whitespace",
			raw:     "  USD , GBP , JPY  ",
			base:    "EUR",
			want:    []string{"USD", "GBP", "JPY"},
			wantLen: 3,
		},
		{
			name:    "lowercase symbols",
			raw:     "usd,gbp",
			base:    "EUR",
			want:    []string{"USD", "GBP"},
			wantLen: 2,
		},
		{
			name:    "excludes base currency from symbols",
			raw:     "USD,EUR,GBP",
			base:    "EUR",
			want:    []string{"USD", "GBP"},
			wantLen: 2,
		},
		{
			name:    "invalid symbols are filtered out",
			raw:     "USD,INVALID,GBP",
			base:    "EUR",
			want:    []string{"USD", "GBP"},
			wantLen: 2,
		},
		{
			name:    "all invalid symbols",
			raw:     "INVALID1,INVALID2",
			base:    "EUR",
			want:    []string{},
			wantLen: 0,
		},
		{
			name:    "symbol equals base currency only",
			raw:     "EUR",
			base:    "EUR",
			want:    []string{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeSymbolsArray(tt.raw, tt.base)
			if len(got) != tt.wantLen {
				t.Errorf("MakeSymbolsArray(%q, %q) returned %d symbols, want %d", tt.raw, tt.base, len(got), tt.wantLen)
			}

			for i, symbol := range tt.want {
				if i >= len(got) || got[i] != symbol {
					t.Errorf("MakeSymbolsArray(%q, %q)[%d] = %q, want %q", tt.raw, tt.base, i, safeGet(got, i), symbol)
				}
			}
		})
	}
}

func safeGet(arr []string, index int) string {
	if index >= len(arr) {
		return "<missing>"
	}
	return arr[index]
}

func TestRatesService_GetLatestRate(t *testing.T) {
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
		name           string
		base           string
		symbols        string
		mockRecord     *ExchangeRateRecord
		mockErr        error
		wantErr        bool
		wantNil        bool
		wantBase       string
		wantRatesCount int
		wantRates      map[string]float64
	}{
		{
			name:           "successful fetch with no symbol filter",
			base:           "EUR",
			symbols:        "",
			mockRecord:     sampleRecord,
			mockErr:        nil,
			wantErr:        false,
			wantNil:        false,
			wantBase:       "EUR",
			wantRatesCount: 4,
		},
		{
			name:       "successful fetch with single symbol filter",
			base:       "EUR",
			symbols:    "USD",
			mockRecord: sampleRecord,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    false,
			wantBase:   "EUR",
			wantRates:  map[string]float64{"USD": 1.08},
		},
		{
			name:       "successful fetch with multiple symbol filter",
			base:       "EUR",
			symbols:    "USD,GBP",
			mockRecord: sampleRecord,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    false,
			wantBase:   "EUR",
			wantRates:  map[string]float64{"USD": 1.08, "GBP": 0.86},
		},
		{
			name:       "filter with non-existent symbol",
			base:       "EUR",
			symbols:    "USD,AUD",
			mockRecord: sampleRecord,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    false,
			wantBase:   "EUR",
			wantRates:  map[string]float64{"USD": 1.08},
		},
		{
			name:       "repository returns nil (not found)",
			base:       "EUR",
			symbols:    "",
			mockRecord: nil,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    true,
		},
		{
			name:       "repository returns error",
			base:       "EUR",
			symbols:    "",
			mockRecord: nil,
			mockErr:    errors.New("database error"),
			wantErr:    true,
			wantNil:    true,
		},
		{
			name:       "invalid base currency defaults to EUR",
			base:       "INVALID",
			symbols:    "",
			mockRecord: sampleRecord,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    false,
			wantBase:   "EUR",
		},
		{
			name:       "lowercase base currency is normalized",
			base:       "eur",
			symbols:    "",
			mockRecord: sampleRecord,
			mockErr:    nil,
			wantErr:    false,
			wantNil:    false,
			wantBase:   "EUR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRatesRepository{
				GetLatestRateFunc: func(base string) (*ExchangeRateRecord, error) {
					return tt.mockRecord, tt.mockErr
				},
			}

			service := NewRatesService(mockRepo)
			got, err := service.GetLatestRate(tt.base, tt.symbols)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil {
				if got != nil {
					t.Errorf("GetLatestRate() = %v, want nil", got)
				}
				return
			}

			if got == nil && !tt.wantNil {
				t.Errorf("GetLatestRate() returned nil, want non-nil")
				return
			}

			if got != nil {
				if tt.wantBase != "" && got.Base != tt.wantBase {
					t.Errorf("GetLatestRate() base = %q, want %q", got.Base, tt.wantBase)
				}

				if tt.wantRatesCount > 0 && len(got.Rates) != tt.wantRatesCount {
					t.Errorf("GetLatestRate() rates count = %d, want %d", len(got.Rates), tt.wantRatesCount)
				}

				if tt.wantRates != nil {
					for symbol, rate := range tt.wantRates {
						if gotRate, ok := got.Rates[symbol]; !ok || gotRate != rate {
							t.Errorf("GetLatestRate() rates[%s] = %v, want %v", symbol, gotRate, rate)
						}
					}
					if len(got.Rates) != len(tt.wantRates) {
						t.Errorf("GetLatestRate() rates count = %d, want %d", len(got.Rates), len(tt.wantRates))
					}
				}
			}
		})
	}
}

func TestCurrenciesContainsExpected(t *testing.T) {
	expectedCurrencies := []string{"EUR", "USD", "GBP", "JPY", "CHF", "AUD", "CAD", "CNY"}

	for _, currency := range expectedCurrencies {
		found := false
		for _, c := range Currencies {
			if c == currency {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Currencies list missing expected currency: %s", currency)
		}
	}
}

func TestCurrenciesNoDuplicates(t *testing.T) {
	seen := make(map[string]bool)
	for _, currency := range Currencies {
		if seen[currency] {
			t.Errorf("Currencies list contains duplicate: %s", currency)
		}
		seen[currency] = true
	}
}
