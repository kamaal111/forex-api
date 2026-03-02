package routers

import (
	"net/http"

	"github.com/kamaal111/forex-api/handlers"
)

func currenciesGroup(mux *http.ServeMux) {
	mux.Handle("/v1/currencies", loggerMiddleware(http.HandlerFunc(handlers.GetCurrencies)))
}
