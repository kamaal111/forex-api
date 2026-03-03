package routers

import (
	"net/http"

	"github.com/kamaal111/forex-api/handlers"
)

func ratesGroup(mux *http.ServeMux) {
	mux.Handle(handlers.LatestPath, loggerMiddleware(http.HandlerFunc(handlers.GetLatest)))
	mux.Handle(handlers.SymbolsPath, loggerMiddleware(http.HandlerFunc(handlers.GetSymbols)))
}
