// Package main is the entry point for the Forex API.
//
// @title           Forex API
// @version         1.0
// @description     API for fetching currency exchange rates.
//
// @host            localhost:8000
// @BasePath        /
package main

import (
	"github.com/kamaal111/forex-api/routers"
)

func main() {
	routers.Start()
}
