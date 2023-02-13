package routers

import (
	"log"
	"net/http"

	"github.com/kamaal111/forex-api/utils"
)

func Start() {
	serverAddress := utils.UnwrapEnvironment("SERVER_ADDRESS")

	mux := http.NewServeMux()
	ratesGroup(mux)
	mux.Handle("/", loggerMiddleware(http.HandlerFunc(notFound)))

	log.Printf("Listening on %s...", serverAddress)

	err := http.ListenAndServe(serverAddress, mux)
	log.Fatal(err)
}
