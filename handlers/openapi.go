package handlers

import (
	"net/http"

	forexdocs "github.com/kamaal111/forex-api/docs"
)

// GetOpenAPISpec serves the OpenAPI specification in YAML format.
//
// @Summary      Download OpenAPI spec
// @Description  Returns the OpenAPI specification for this API in YAML format.
// @Tags         openapi
// @Produce      application/x-yaml
// @Success      200  {string}  string  "OpenAPI spec in YAML format"
// @Router       /openapi.yaml [get]
func GetOpenAPISpec(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/x-yaml")
	writer.Write(forexdocs.SwaggerYAML)
}
