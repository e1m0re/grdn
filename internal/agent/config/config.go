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
	Key            string
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
	flag.StringVar(&config.Key, "k", "", "key to use for encryption")
	flag.Parse()

	if envServerAddr := os.Getenv("ADDRESS"); envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			reportInterval = uint(envValue)
		}
	}
	config.ReportInterval = time.Duration(reportInterval) * time.Second

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			pollInterval = uint(envValue)
		}
	}
	config.PollInterval = time.Duration(pollInterval) * time.Second

	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

	return &config
}
