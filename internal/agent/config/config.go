// Package config contains instruments for configuration of clients application.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultServerAddr     = "localhost:8080"
	defaultReportInterval = 10
	defaultPollInterval   = 2
	defaultKey            = ""
	defaultRateLimit      = 1
	defaultPublicKey      = ""

	envConfigFileName     = "CONFIG"
	envServerAddrName     = "ADDRESS"
	envReportIntervalName = "REPORT_INTERVAL"
	envPollInterval       = "POLL_INTERVAL"
	envKeyName            = "KEY"
	envRateLimit          = "RATE_LIMIT"
	envCryptoKeyName      = "CRYPTO_KEY"
)

type Config struct {
	Key            string
	PublicKeyFile  string        `json:"crypto_key"`
	ServerAddr     string        `json:"address"`
	PollInterval   time.Duration `json:"poll_interval"`
	ReportInterval time.Duration `json:"report_interval"`
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
	if envConfigFile := os.Getenv(envConfigFileName); envConfigFile != "" {
		configFile = envConfigFile
	}
	if configFile != "" {
		err := updateConfigFromFile(&config, configFile)
		if err != nil {
			return nil, err
		}
	}

	flag.StringVar(&config.ServerAddr, "a", defaultServerAddr, "address and port to run server")
	flag.UintVar(&reportInterval, "r", defaultReportInterval, "frequency of sending metrics to the server")
	flag.UintVar(&pollInterval, "p", defaultPollInterval, "frequency of polling metrics from the package")
	flag.StringVar(&config.Key, "k", defaultKey, "key to use for encryption")
	flag.IntVar(&config.RateLimit, "l", defaultRateLimit, "limit of threads count")
	flag.StringVar(&config.PublicKeyFile, "crypto-key", defaultPublicKey, "public key file path")
	flag.Parse()

	if envServerAddr := os.Getenv(envServerAddrName); envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	if envReportInterval := os.Getenv(envReportIntervalName); envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			reportInterval = uint(envValue)
		}
	}
	config.ReportInterval = time.Duration(reportInterval) * time.Second

	if envPollInterval := os.Getenv(envPollInterval); envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			pollInterval = uint(envValue)
		}
	}
	config.PollInterval = time.Duration(pollInterval) * time.Second

	if envKey := os.Getenv(envKeyName); envKey != "" {
		config.Key = envKey
	}

	if envRateLimit := os.Getenv(envRateLimit); envRateLimit != "" {
		envValue, err := strconv.Atoi(envRateLimit)
		if err == nil {
			config.RateLimit = envValue
		}
	}
	if config.RateLimit <= 0 {
		config.RateLimit = 1
	}

	if envCryptoKey := os.Getenv(envCryptoKeyName); envCryptoKey != "" {
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
