package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/e1m0re/grdn/internal/gvar"
	"github.com/e1m0re/grdn/internal/server"
	"github.com/e1m0re/grdn/internal/server/config"
	"github.com/e1m0re/grdn/internal/signals"
	"github.com/e1m0re/grdn/internal/storage"
	"github.com/e1m0re/grdn/internal/storage/store"
)

func main() {
	gvar.PrintWelcome()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go signals.HandleOSSignals(cancel)

	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error("error init configuration", slog.String("error", err.Error()))
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	s, err := initializeStore(ctx, cfg)
	if err != nil {
		slog.Error("error init store", slog.String("error", err.Error()))
		return
	}

	srv := server.NewServer(cfg, s)
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

func initializeStore(ctx context.Context, cfg *config.Config) (store.Store, error) {
	storeType := storage.TypeMemory
	path := cfg.FileStoragePath
	if len(cfg.DatabaseDSN) > 0 {
		storeType = storage.TypePostgres
		path = cfg.DatabaseDSN
	}
	newStore, err := store.NewStore(ctx, &storage.Config{
		Path:     path,
		Type:     storeType,
		Interval: cfg.StoreInternal,
	})
	if err != nil {
		return nil, err
	}

	return newStore, err
}
