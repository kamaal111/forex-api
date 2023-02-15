package utils

import (
	"log"
	"os"
	"strings"
)

func UnwrapEnvironment(keys ...string) string {
	for _, key := range keys {
		value := os.Getenv(key)
		if value != "" {
			return value
		}
	}

	log.Fatalf("%s not defined in environment\n", strings.Join(keys, ", "))
	return "" // unreachable code
}
