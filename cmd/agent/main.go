package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/e1m0re/grdn/internal/agent/app"
	"github.com/e1m0re/grdn/internal/agent/config"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		slog.Error(err.Error())
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("failed to sync logger: %s", err.Error())
		}
	}(logger)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cfg := config.GetConfig()

	app1 := app.NewApp(cfg)

	err := app1.Start(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
