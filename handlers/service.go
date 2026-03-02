package handlers

import (
	"strings"

	"github.com/kamaal111/forex-api/utils"
)

type ExchangeRateRecord struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type SymbolsRecord struct {
	Date    string   `json:"date" firestore:"date"`
	Symbols []string `json:"symbols" firestore:"symbols"`
}

type NamedSymbol struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type NamedSymbolsRecord struct {
	Date    string        `json:"date"`
	Symbols []NamedSymbol `json:"symbols"`
}

type RatesRepository interface {
	GetLatestRate(base string) (*ExchangeRateRecord, error)
	GetAllSymbols() (*SymbolsRecord, error)
}

type RatesService struct {
	Repository RatesRepository
}

func NewRatesService(repo RatesRepository) *RatesService {
	return &RatesService{Repository: repo}
}

func (s *RatesService) GetLatestRate(base string, symbols string) (*ExchangeRateRecord, error) {
	normalizedBase := NormalizeBase(base)

	record, err := s.Repository.GetLatestRate(normalizedBase)
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, nil
	}

	symbolsArray := MakeSymbolsArray(symbols, normalizedBase)
	if len(symbolsArray) > 0 {
		filteredRecord := &ExchangeRateRecord{
			Base:  record.Base,
			Date:  record.Date,
			Rates: make(map[string]float64),
		}
		for _, symbol := range symbolsArray {
			if rate, ok := record.Rates[symbol]; ok {
				filteredRecord.Rates[symbol] = rate
			}
		}
		return filteredRecord, nil
	}

	return record, nil
}

func (s *RatesService) GetAllSymbols() (*SymbolsRecord, error) {
	return s.Repository.GetAllSymbols()
}

func (s *RatesService) GetAllNamedSymbols() (*NamedSymbolsRecord, error) {
	record, err := s.Repository.GetAllSymbols()
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}

	named := make([]NamedSymbol, 0, len(record.Symbols))
	for _, symbol := range record.Symbols {
		name, ok := CurrencyNames[symbol]
		if !ok {
			name = symbol
		}
		named = append(named, NamedSymbol{Symbol: symbol, Name: name})
	}

	return &NamedSymbolsRecord{Date: record.Date, Symbols: named}, nil
}

func NormalizeBase(base string) string {
	normalized := strings.ToUpper(strings.TrimSpace(base))
	if !utils.ArrayContains(Currencies, normalized) {
		return "EUR"
	}
	return normalized
}

func MakeSymbolsArray(raw string, base string) []string {
	symbols := strings.ToUpper(strings.TrimSpace(raw))
	if len(symbols) == 0 || symbols == "*" {
		return []string{}
	}

	var symbolsArray []string
	for item := range strings.SplitSeq(symbols, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != base && utils.ArrayContains(Currencies, trimmed) {
			symbolsArray = append(symbolsArray, trimmed)
		}
	}
	return symbolsArray
}

var Currencies = []string{
	"EUR",
	"USD",
	"JPY",
	"BGN",
	"CYP",
	"CZK",
	"DKK",
	"EEK",
	"GBP",
	"HUF",
	"LTL",
	"LVL",
	"MTL",
	"PLN",
	"ROL",
	"RON",
	"SEK",
	"SIT",
	"SKK",
	"CHF",
	"ISK",
	"ILS",
	"NOK",
	"HRK",
	"RUB",
	"TRL",
	"TRY",
	"AUD",
	"BRL",
	"CAD",
	"CNY",
	"HKD",
	"IDR",
	"INR",
	"KRW",
	"MXN",
	"MYR",
	"NZD",
	"PHP",
	"SGD",
	"THB",
	"ZAR",
}

var CurrencyNames = map[string]string{
	"EUR": "Euro",
	"USD": "US Dollar",
	"JPY": "Japanese Yen",
	"BGN": "Bulgarian Lev",
	"CYP": "Cypriot Pound",
	"CZK": "Czech Koruna",
	"DKK": "Danish Krone",
	"EEK": "Estonian Kroon",
	"GBP": "British Pound Sterling",
	"HUF": "Hungarian Forint",
	"LTL": "Lithuanian Litas",
	"LVL": "Latvian Lats",
	"MTL": "Maltese Lira",
	"PLN": "Polish Zloty",
	"ROL": "Romanian Leu (old)",
	"RON": "Romanian Leu",
	"SEK": "Swedish Krona",
	"SIT": "Slovenian Tolar",
	"SKK": "Slovak Koruna",
	"CHF": "Swiss Franc",
	"ISK": "Icelandic Krona",
	"ILS": "Israeli New Shekel",
	"NOK": "Norwegian Krone",
	"HRK": "Croatian Kuna",
	"RUB": "Russian Ruble",
	"TRL": "Turkish Lira (old)",
	"TRY": "Turkish Lira",
	"AUD": "Australian Dollar",
	"BRL": "Brazilian Real",
	"CAD": "Canadian Dollar",
	"CNY": "Chinese Yuan",
	"HKD": "Hong Kong Dollar",
	"IDR": "Indonesian Rupiah",
	"INR": "Indian Rupee",
	"KRW": "South Korean Won",
	"MXN": "Mexican Peso",
	"MYR": "Malaysian Ringgit",
	"NZD": "New Zealand Dollar",
	"PHP": "Philippine Peso",
	"SGD": "Singapore Dollar",
	"THB": "Thai Baht",
	"ZAR": "South African Rand",
}
