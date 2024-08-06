// Package config define server configuration.
package config

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/netip"
	"os"
	"time"
)

const (
	defaultServerAddr      = "localhost:8080"
	defaultVerboseMode     = false
	defaultStoreInternal   = 300 * time.Second
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestoreData     = true
	defaultDatabaseDSN     = ""
	defaultKey             = ""
	defaultPrivateKeyFile  = ""
	defaultTrustedSubnet   = ""

	envConfigFileName      = "CONFIG"
	envRunAddrName         = "ADDRESS"
	envStoreIntervalName   = "STORE_INTERVAL"
	envFileStoragePathName = "FILE_STORAGE_PATH"
	envRestoreDataName     = "RESTORE"
	envDatabaseDSNName     = "DATABASE_DSN"
	envKeyName             = "KEY"
	envCryptoKeyName       = "CRYPTO_KEY"
	envTrustedSubnet       = "TRUSTED_SUBNET"
)

type Config struct {
	FileStoragePath string `json:"store_file"`
	LoggerLevel     string
	ServerAddr      string `json:"address"`
	DatabaseDSN     string `json:"database_dsn"`
	Key             string
	PrivateKeyFile  string        `json:"crypto_key"`
	TrustedSubnet   string        `json:"trusted_subnet"`
	StoreInternal   time.Duration `json:"store_interval"`
	LogLevel        slog.Level
	RestoreData     bool `json:"restore"`
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
	flag.BoolVar(&config.VerboseMode, "v", defaultVerboseMode, "Torn on extended logging mode")
	flag.DurationVar(&config.StoreInternal, "i", defaultStoreInternal, "time interval to save data to HDD")
	flag.StringVar(&config.FileStoragePath, "f", defaultFileStoragePath, "file path for DB file")
	flag.BoolVar(&config.RestoreData, "r", defaultRestoreData, "save or don't save data to HDD on shutdown")
	flag.StringVar(&config.DatabaseDSN, "d", defaultDatabaseDSN, "database connection string")
	flag.StringVar(&config.Key, "k", defaultKey, "key to use for encryption")
	flag.StringVar(&config.PrivateKeyFile, "crypto-key", defaultPrivateKeyFile, "public key file path")
	flag.StringVar(&config.TrustedSubnet, "t", defaultTrustedSubnet, "trusted subnet of clients")
	flag.Parse()

	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		config.ServerAddr = envRunAddr
	}

	if config.VerboseMode {
		config.LogLevel = slog.LevelDebug
		config.LoggerLevel = "debug"
	}

	if envStoreInterval := os.Getenv(envStoreIntervalName); envStoreInterval != "" {
		value, err := time.ParseDuration(envStoreInterval)
		if err == nil {
			config.StoreInternal = value
		}
	}

	if envFileStoragePath := os.Getenv(envFileStoragePathName); envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	if envRestoreData := os.Getenv(envRestoreDataName); envRestoreData != "" {
		config.RestoreData = envRestoreData == "true"
	}

	if envDatabaseDSN := os.Getenv(envDatabaseDSNName); envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	if envKey := os.Getenv(envKeyName); envKey != "" {
		config.Key = envKey
	}

	if envCryptoKey := os.Getenv(envCryptoKeyName); envCryptoKey != "" {
		config.PrivateKeyFile = envCryptoKey
	}

	if envTrustedSubnet := os.Getenv(envTrustedSubnet); envTrustedSubnet != "" {
		config.TrustedSubnet = envTrustedSubnet
	}
	if _, err := netip.ParsePrefix(config.TrustedSubnet); err != nil {
		return nil, err
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
