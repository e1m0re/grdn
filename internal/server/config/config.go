package config

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
}

func GetConfig() *Config {
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
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.ServerAddr = envRunAddr
	}

	if config.VerboseMode {
		config.LogLevel = slog.LevelDebug
		config.LoggerLevel = "debug"
	}

	envStoreInterval := os.Getenv("STORE_INTERVAL")
	if envStoreInterval != "" {
		value, err := time.ParseDuration(envStoreInterval)
		if err == nil {
			config.StoreInternal = value
		}
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	envRestoreData := os.Getenv("RESTORE")
	if envRestoreData != "" {
		config.RestoreData = envRestoreData == "true"
	}

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	return &config
}
