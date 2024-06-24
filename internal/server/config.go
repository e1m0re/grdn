package server

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	FileStoragePath string
	LogLevel        slog.Level
	LoggerLevel     string
	RestoreData     bool
	ServerAddr      string
	StoreInternal   time.Duration
	VerboseMode     bool
	DatabaseDSN     string
	Key             string
}

func InitConfig() *Config {
	config := Config{
		LogLevel:    slog.LevelInfo,
		LoggerLevel: "info",
	}

	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.BoolVar(&config.VerboseMode, "v", false, "Torn on extended logging mode")
	flag.DurationVar(&config.StoreInternal, "i", 300*time.Second, "time interval to save data to HDD")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "file path for DB file")
	flag.BoolVar(&config.RestoreData, "r", true, "save or don't save data to HDD on shutdown")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&config.Key, "k", "", "key to use for encryption")
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

	return &config
}
