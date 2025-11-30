package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/kamaal111/forex-api/database"
	"github.com/kamaal111/forex-api/utils"
)

type exchangeRateRecord struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func GetLatest(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	client, err := database.CreateClient(ctx)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	base := strings.ToUpper(strings.TrimSpace(request.URL.Query().Get("base")))
	if !utils.ArrayContains(CURRENCIES, base) {
		base = "EUR"
	}

	documents := client.Collection("exchange_rates").OrderBy("date", firestore.Desc).Where("base", "==", base).Limit(1).Documents(ctx)
	var document *firestore.DocumentSnapshot
	for document == nil {
		document, err = documents.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if document == nil {
		utils.ErrorHandler(writer, "Rates not found", http.StatusNotFound)
		return
	}

	var record exchangeRateRecord
	err = document.DataTo(&record)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
	}

	symbols := makeSymbolsArray(request.URL.Query().Get("symbols"), base)
	if len(symbols) > 0 {
		recordCopy := record
		recordCopy.Rates = make(map[string]float64)
		for _, symbol := range symbols {
			recordCopy.Rates[symbol] = record.Rates[symbol]
		}
		record = recordCopy
	}

	output, err := json.Marshal(record)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "application/json")
	writer.Write(output)
}

func makeSymbolsArray(raw string, base string) []string {
	symbols := strings.ToUpper(strings.TrimSpace(raw))
	if len(symbols) == 0 {
		return []string{}
	}

	var symbolsArray []string
	for item := range strings.SplitSeq(symbols, ",") {
		if item != base && utils.ArrayContains(CURRENCIES, item) {
			symbolsArray = append(symbolsArray, item)
		}
	}
	return symbolsArray
}

var CURRENCIES = []string{
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
	"ILS",
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
