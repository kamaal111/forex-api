package routers

import (
	"net/http"

	"github.com/kamaal111/forex-api/utils"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	utils.ErrorHandler(w, "Not found", http.StatusNotFound)
}
