package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr     string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func GetConfig() *Config {
	config := Config{}

	var (
		pollInterval   uint
		reportInterval uint
	)
	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.UintVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.UintVar(&pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()

	envServerAddr := os.Getenv("ADDRESS")
	if envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	envReportInterval := os.Getenv("REPORT_INTERVAL")
	if envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			reportInterval = uint(envValue)
		}
	}
	config.ReportInterval = time.Duration(reportInterval) * time.Second

	envPollInterval := os.Getenv("POLL_INTERVAL")
	if envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			pollInterval = uint(envValue)
		}
	}
	config.PollInterval = time.Duration(pollInterval) * time.Second

	return &config
}
