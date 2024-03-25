package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/http-server/server"
	"github.com/e1m0re/grdn/internal/logger"
	"github.com/e1m0re/grdn/internal/storage"
)

func main() {
	parameters := config()

	mySlogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parameters.logLevel}))
	slog.SetDefault(mySlogger)

	if err := logger.Initialize(parameters.loggerLevel); err != nil {
		slog.Error(err.Error())
		return
	}

	store := storage.NewMemStorage(parameters.storeInternal == 0, parameters.fileStoragePath)
	if parameters.restoreData {
		err := store.LoadStorageFromFile()
		if err != nil {
			slog.Error(fmt.Sprintf("error restore data: %s", err))
		}
	}

	httpServer := server.NewServer(parameters.serverAddr, store)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.Info(fmt.Sprintf("Running server on %s", parameters.serverAddr))
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			return err
		}

		return store.DumpStorageToFile()
	})

	if parameters.storeInternal > 0 {
		g.Go(func() error {
			for {
				select {
				case <-gCtx.Done():
					return nil
				case <-time.After(parameters.storeInternal):
					err := store.DumpStorageToFile()
					if err != nil {
						slog.Error(fmt.Sprintf("error autosave: %s", err))
					}
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info(err.Error())
			return
		}

		slog.Error(err.Error())
	}
}
