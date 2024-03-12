package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {
	defaultRunAddr := "localhost:8080"
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		defaultRunAddr = envRunAddr
	}

	flag.StringVar(&flagRunAddr, "a", defaultRunAddr, "address and port to run server")
	flag.Parse()
}
