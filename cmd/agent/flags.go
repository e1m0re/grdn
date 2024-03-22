package main

import (
	"flag"
	"os"
	"strconv"
)

type parameters struct {
	serverAddr     string
	reportInterval uint
	pollInterval   uint
}

func config() *parameters {
	options := parameters{
		serverAddr:     "",
		reportInterval: 0,
		pollInterval:   0,
	}

	flag.StringVar(&options.serverAddr, "a", "localhost:8080", "address and port to run server")
	flag.UintVar(&options.reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.UintVar(&options.pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()

	envServerAddr := os.Getenv("ADDRESS")
	if envServerAddr != "" {
		options.serverAddr = envServerAddr
	}

	envReportInterval := os.Getenv("REPORT_INTERVAL")
	if envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			options.reportInterval = uint(envValue)
		}
	}

	envPollInterval := os.Getenv("POLL_INTERVAL")
	if envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			options.pollInterval = uint(envValue)
		}
	}

	return &options
}
