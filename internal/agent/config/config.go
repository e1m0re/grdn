// Package config contains instruments for configuration of clients application.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Key            string
	PublicKeyFile  string        `yaml:"crypto_key"`
	ServerAddr     string        `yaml:"address"`
	PollInterval   time.Duration `yaml:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval"`
	RateLimit      int
}

// InitConfig initializes the clients application configuration.
func InitConfig() (*Config, error) {
	config := Config{}

	var (
		pollInterval   uint
		reportInterval uint
	)

	var configFile string
	flag.StringVar(&configFile, "c", "", "config file (JSON)")
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
	}
	if configFile != "" {
		err := updateConfigFromFile(&config, configFile)
		if err != nil {
			return nil, err
		}
	}

	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.UintVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.UintVar(&pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.StringVar(&config.Key, "k", "", "key to use for encryption")
	flag.IntVar(&config.RateLimit, "l", 1, "limit of threads count")
	flag.StringVar(&config.PublicKeyFile, "crypto-key", "", "public key file path")
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

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		envValue, err := strconv.Atoi(envRateLimit)
		if err == nil {
			config.RateLimit = envValue
		}
	}
	if config.RateLimit <= 0 {
		config.RateLimit = 1
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		config.PublicKeyFile = envCryptoKey
	}

	return &config, nil
}

func updateConfigFromFile(c *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	return decoder.Decode(c)
}
