package main

import (
	"flag"
	"os"
)

type parameters struct {
	serverAddr  string
	verboseMode bool
	loggerLevel string
}

func config() *parameters {
	options := parameters{
		serverAddr:  "",
		verboseMode: false,
		loggerLevel: "info",
	}

	flag.StringVar(&options.serverAddr, "a", "localhost:8080", "address and port to run server")
	flag.BoolVar(&options.verboseMode, "v", false, "Torn on extended logging mode")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		options.serverAddr = envRunAddr
	}

	if options.verboseMode {
		options.loggerLevel = "debug"
	}

	return &options
}
