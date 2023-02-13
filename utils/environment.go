package utils

import (
	"log"
	"os"
)

func UnwrapEnvironment(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not defined in environment\n", key)
	}

	return value
}
