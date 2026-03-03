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

type CurrencyInfo struct {
	Name string
	Sign string
}

type NamedSymbol struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Sign   string `json:"sign"`
}

type CurrenciesRecord struct {
	Date string        `json:"date"`
	Data []NamedSymbol `json:"data"`
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

func (s *RatesService) GetAllNamedSymbols() (*CurrenciesRecord, error) {
	record, err := s.Repository.GetAllSymbols()
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}

	named := make([]NamedSymbol, 0, len(record.Symbols))
	for _, symbol := range record.Symbols {
		if info, ok := CurrencyNames[symbol]; ok {
			named = append(named, NamedSymbol{Symbol: symbol, Name: info.Name, Sign: info.Sign})
		}
	}

	return &CurrenciesRecord{Date: record.Date, Data: named}, nil
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

var CurrencyNames = map[string]CurrencyInfo{
	"EUR": {Name: "Euro", Sign: "€"},
	"USD": {Name: "US Dollar", Sign: "$"},
	"JPY": {Name: "Japanese Yen", Sign: "¥"},
	"BGN": {Name: "Bulgarian Lev", Sign: "лв"},
	"CYP": {Name: "Cypriot Pound", Sign: "£"},
	"CZK": {Name: "Czech Koruna", Sign: "Kč"},
	"DKK": {Name: "Danish Krone", Sign: "kr"},
	"EEK": {Name: "Estonian Kroon", Sign: "kr"},
	"GBP": {Name: "British Pound Sterling", Sign: "£"},
	"HUF": {Name: "Hungarian Forint", Sign: "Ft"},
	"LTL": {Name: "Lithuanian Litas", Sign: "Lt"},
	"LVL": {Name: "Latvian Lats", Sign: "Ls"},
	"MTL": {Name: "Maltese Lira", Sign: "₤"},
	"PLN": {Name: "Polish Zloty", Sign: "zł"},
	"ROL": {Name: "Romanian Leu (old)", Sign: "lei"},
	"RON": {Name: "Romanian Leu", Sign: "lei"},
	"SEK": {Name: "Swedish Krona", Sign: "kr"},
	"SIT": {Name: "Slovenian Tolar", Sign: "SIT"},
	"SKK": {Name: "Slovak Koruna", Sign: "Sk"},
	"CHF": {Name: "Swiss Franc", Sign: "Fr"},
	"ISK": {Name: "Icelandic Krona", Sign: "kr"},
	"ILS": {Name: "Israeli New Shekel", Sign: "₪"},
	"NOK": {Name: "Norwegian Krone", Sign: "kr"},
	"HRK": {Name: "Croatian Kuna", Sign: "kn"},
	"RUB": {Name: "Russian Ruble", Sign: "₽"},
	"TRL": {Name: "Turkish Lira (old)", Sign: "₤"},
	"TRY": {Name: "Turkish Lira", Sign: "₺"},
	"AUD": {Name: "Australian Dollar", Sign: "$"},
	"BRL": {Name: "Brazilian Real", Sign: "R$"},
	"CAD": {Name: "Canadian Dollar", Sign: "$"},
	"CNY": {Name: "Chinese Yuan", Sign: "¥"},
	"HKD": {Name: "Hong Kong Dollar", Sign: "$"},
	"IDR": {Name: "Indonesian Rupiah", Sign: "Rp"},
	"INR": {Name: "Indian Rupee", Sign: "₹"},
	"KRW": {Name: "South Korean Won", Sign: "₩"},
	"MXN": {Name: "Mexican Peso", Sign: "$"},
	"MYR": {Name: "Malaysian Ringgit", Sign: "RM"},
	"NZD": {Name: "New Zealand Dollar", Sign: "$"},
	"PHP": {Name: "Philippine Peso", Sign: "₱"},
	"SGD": {Name: "Singapore Dollar", Sign: "$"},
	"THB": {Name: "Thai Baht", Sign: "฿"},
	"ZAR": {Name: "South African Rand", Sign: "R"},
}

var Currencies = func() []string {
	keys := make([]string, 0, len(CurrencyNames))
	for k := range CurrencyNames {
		keys = append(keys, k)
	}
	return keys
}()
