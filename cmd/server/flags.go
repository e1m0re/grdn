package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {
	var defaultA = "localhost:8080"
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		defaultA = envRunAddr
	}

	flag.StringVar(&flagRunAddr, "a", defaultA, "address and port to run server")
	flag.Parse()
}
