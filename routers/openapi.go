package routers

import (
	"net/http"

	"github.com/kamaal111/forex-api/handlers"
)

func openapiGroup(mux *http.ServeMux) {
	mux.Handle(handlers.OpenAPISpecPath, loggerMiddleware(http.HandlerFunc(handlers.GetOpenAPISpec)))
}
