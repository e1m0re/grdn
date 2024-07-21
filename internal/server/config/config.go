// Package config define server configuration.
package config

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	FileStoragePath string `yaml:"store_file"`
	LoggerLevel     string
	ServerAddr      string `yaml:"address"`
	DatabaseDSN     string `yaml:"database_dsn"`
	Key             string
	PrivateKeyFile  string        `yaml:"crypto_key"`
	StoreInternal   time.Duration `yaml:"store_interval"`
	LogLevel        slog.Level
	RestoreData     bool `yaml:"restore"`
	VerboseMode     bool
}

// InitConfig initializes the server configuration.
func InitConfig() (*Config, error) {
	config := Config{
		LogLevel:    slog.LevelInfo,
		LoggerLevel: "info",
	}

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
	flag.BoolVar(&config.VerboseMode, "v", false, "Torn on extended logging mode")
	flag.DurationVar(&config.StoreInternal, "i", 300*time.Second, "time interval to save data to HDD")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "file path for DB file")
	flag.BoolVar(&config.RestoreData, "r", true, "save or don't save data to HDD on shutdown")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&config.Key, "k", "", "key to use for encryption")
	flag.StringVar(&config.PrivateKeyFile, "crypto-key", "", "public key file path")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.ServerAddr = envRunAddr
	}

	if config.VerboseMode {
		config.LogLevel = slog.LevelDebug
		config.LoggerLevel = "debug"
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := time.ParseDuration(envStoreInterval)
		if err == nil {
			config.StoreInternal = value
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	if envRestoreData := os.Getenv("RESTORE"); envRestoreData != "" {
		config.RestoreData = envRestoreData == "true"
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		config.PrivateKeyFile = envCryptoKey
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
