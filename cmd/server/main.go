package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/e1m0re/grdn/internal/gvar"
	"github.com/e1m0re/grdn/internal/server"
	"github.com/e1m0re/grdn/internal/server/config"
	"github.com/e1m0re/grdn/internal/server/storage"
	"github.com/e1m0re/grdn/internal/server/storage/store"
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

	cfg := config.InitConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	initializeStore(cfg)

	srv := server.NewServer(cfg)
	err := srv.Start(ctx)
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

func initializeStore(cfg *config.Config) {
	storeType := storage.TypeMemory
	path := cfg.FileStoragePath
	if len(cfg.DatabaseDSN) > 0 {
		storeType = storage.TypePostgres
		path = cfg.DatabaseDSN
	}
	err := store.Initialize(&storage.Config{
		Path:     path,
		Type:     storeType,
		Interval: cfg.StoreInternal,
	})
	if err != nil {
		panic(err)
	}
}
