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

type RatesRepository interface {
	GetLatestRate(base string) (*ExchangeRateRecord, error)
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

func NormalizeBase(base string) string {
	normalized := strings.ToUpper(strings.TrimSpace(base))
	if !utils.ArrayContains(Currencies, normalized) {
		return "EUR"
	}
	return normalized
}

func MakeSymbolsArray(raw string, base string) []string {
	symbols := strings.ToUpper(strings.TrimSpace(raw))
	if len(symbols) == 0 {
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
