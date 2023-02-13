package routers

import (
	"net/http"

	"github.com/kamaal111/forex-api/handlers"
)

func ratesGroup(mux *http.ServeMux) {
	mux.Handle("/v1/rates/latest", loggerMiddleware(http.HandlerFunc(handlers.GetLatest)))
}
