package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

type parameters struct {
	fileStoragePath string
	logLevel        slog.Level
	loggerLevel     string
	restoreData     bool
	serverAddr      string
	storeInternal   time.Duration
	verboseMode     bool
}

func config() *parameters {
	options := parameters{
		logLevel:    slog.LevelInfo,
		loggerLevel: "info",
	}

	flag.StringVar(&options.serverAddr, "a", "localhost:8080", "address and port to run server")
	flag.BoolVar(&options.verboseMode, "v", false, "Torn on extended logging mode")
	flag.DurationVar(&options.storeInternal, "i", 300*time.Second, "time interval to save data to HDD")
	flag.StringVar(&options.fileStoragePath, "f", "/tmp/metrics-db.json", "file path for DB file")
	flag.BoolVar(&options.restoreData, "r", true, "save or don't save data to HDD on shutdown")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		options.serverAddr = envRunAddr
	}

	if options.verboseMode {
		options.logLevel = slog.LevelDebug
		options.loggerLevel = "debug"
	}

	envStoreInterval := os.Getenv("STORE_INTERVAL")
	if envStoreInterval != "" {
		value, err := time.ParseDuration(envStoreInterval)
		if err == nil {
			options.storeInternal = value
		}
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		options.fileStoragePath = envFileStoragePath
	}

	envRestoreData := os.Getenv("RESTORE")
	if envRestoreData != "" {
		options.restoreData = envRestoreData == "true"
	}

	return &options
}
