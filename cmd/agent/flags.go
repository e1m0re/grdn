package main

import (
	"flag"
	"os"
	"strconv"
)

type Options struct {
	flagRunAddr    string
	reportInterval uint
	pollInterval   uint
}

var agentOptions = new(Options)

func parseFlags() {
	var defaultA = "localhost:8080"
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		defaultA = envRunAddr
	}
	flag.StringVar(&agentOptions.flagRunAddr, "a", defaultA, "address and port to run server")

	var defaultR = uint(10)
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			defaultR = uint(envValue)
		}
	}
	flag.UintVar(&agentOptions.reportInterval, "r", defaultR, "frequency of sending metrics to the server")

	var defaultP = uint(2)
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			defaultP = uint(envValue)
		}
	}
	flag.UintVar(&agentOptions.pollInterval, "p", defaultP, "frequency of polling metrics from the package")

	flag.Parse()
}
