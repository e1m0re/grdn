package main

import (
	"context"
	"errors"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/e1m0re/grdn/internal/logger"
	"github.com/e1m0re/grdn/internal/server/config"
	"github.com/e1m0re/grdn/internal/server/server"
)

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

	mySlogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(mySlogger)
	if err := logger.Initialize(cfg.LoggerLevel); err != nil {
		slog.Error(err.Error())
		return
	}

	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		slog.Error(err.Error())
	}

	err = srv.Start(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info(err.Error())
		return
	}
	slog.Error(err.Error())
}
