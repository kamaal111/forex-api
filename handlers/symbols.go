package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kamaal111/forex-api/utils"
)

func GetSymbols(writer http.ResponseWriter, request *http.Request) {
	output, err := json.Marshal(Currencies)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "application/json")
	writer.Write(output)
}
