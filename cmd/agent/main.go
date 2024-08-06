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
	"github.com/e1m0re/grdn/internal/gvar"
	"github.com/e1m0re/grdn/internal/service"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	gvar.PrintWelcome()

	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	services, err := service.NewAgentServices(cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	app1, err := app.NewApp(cfg, services)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = app1.Start(ctx, "grpc"); err != nil {
		slog.Error(err.Error())
	}
}
