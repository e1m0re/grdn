package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.uber.org/zap"

	"github.com/e1m0re/grdn/internal/agent/app"
	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/gvar"
	"github.com/e1m0re/grdn/internal/signals"
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

	mySlogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(mySlogger)
}

func main() {
	gvar.PrintWelcome()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go signals.HandleOSSignals(cancel)

	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	app1, err := app.NewApp(cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = app1.Start(ctx); err != nil {
		slog.Error(err.Error())
	}
}
