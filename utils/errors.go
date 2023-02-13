package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func ErrorHandler(w http.ResponseWriter, message string, code int) {
	errorResponse := Error{
		Message: message,
		Status:  code,
	}
	log.Printf("failure message: %s; code: %d", message, code)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(errorResponse.Status)
	json.NewEncoder(w).Encode(errorResponse)
}
