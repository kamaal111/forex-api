package routers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kamaal111/forex-api/utils"
)

func Start() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = fmt.Sprintf(":%s", utils.UnwrapEnvironment("PORT"))
	}

	mux := http.NewServeMux()
	ratesGroup(mux)
	mux.Handle("/", loggerMiddleware(http.HandlerFunc(notFound)))

	log.Printf("Listening on %s...", serverAddress)

	err := http.ListenAndServe(serverAddress, mux)
	log.Fatal(err)
}
