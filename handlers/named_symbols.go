package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kamaal111/forex-api/database"
	"github.com/kamaal111/forex-api/utils"
)

// GetCurrencies handles requests for currencies with names and signs.
//
// @Summary      Get currencies with names and signs
// @Description  Returns all available currencies with their human-readable names and currency signs.
// @Tags         currencies
// @Produce      json
// @Success      200  {object}  CurrenciesRecord
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /v1/currencies [get]
func GetCurrencies(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	client, err := database.CreateClient(ctx)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	repo := NewFirestoreRatesRepository(ctx, client)
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
