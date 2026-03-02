package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kamaal111/forex-api/database"
	"github.com/kamaal111/forex-api/utils"
)

func GetNamedSymbols(writer http.ResponseWriter, request *http.Request) {
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
