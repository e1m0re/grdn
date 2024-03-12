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

var agentOptions = &Options{}

func parseFlags() {
	defaultRunAddr := "localhost:8080"
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		defaultRunAddr = envRunAddr
	}

	flag.StringVar(&agentOptions.flagRunAddr, "a", defaultRunAddr, "address and port to run server")

	defaultReportInterval := uint(10)

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			defaultReportInterval = uint(envValue)
		}
	}

	flag.UintVar(&agentOptions.reportInterval, "r", defaultReportInterval, "frequency of sending metrics to the server")

	defaultPollInterval := uint(2)

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			defaultPollInterval = uint(envValue)
		}
	}

	flag.UintVar(&agentOptions.pollInterval, "p", defaultPollInterval, "frequency of polling metrics from the package")

	flag.Parse()
}
