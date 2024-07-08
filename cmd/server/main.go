package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/e1m0re/grdn/internal/gvar"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/e1m0re/grdn/internal/server"
)

func main() {
	gvar.PrintWelcome()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cfg := server.InitConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	srv, err := server.NewServer(cfg)
	if err != nil {
		slog.Error(err.Error())
	}

	err = srv.Start(ctx)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info(err.Error())
			return
		}

		slog.Error("error",
			slog.String("error", fmt.Sprintf("%v", err)),
			slog.String("stack", string(debug.Stack())),
		)
	}
}
