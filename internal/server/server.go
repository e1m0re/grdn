package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/db/migrations"
	appHandler "github.com/e1m0re/grdn/internal/server/handler"
	"github.com/e1m0re/grdn/internal/service"
	"github.com/e1m0re/grdn/internal/storage"
	"github.com/e1m0re/grdn/internal/utils"
)

type Server struct {
	cfg        *Config
	httpServer *http.Server
	store      storage.Store
}

func NewServer(cfg *Config) (*Server, error) {
	srv := &Server{
		cfg: cfg,
	}

	err := srv.initStore()

	services := service.NewServices(srv.store)
	handler := appHandler.NewHandler(services)
	srv.httpServer = &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: handler.NewRouter(cfg.Key),
	}

	return srv, err
}

func (srv *Server) initStore() error {
	if srv.cfg.DatabaseDSN != "" {
		err := srv.migrate()
		if err != nil {
			return err
		}

		store, err := storage.NewDBStorage(srv.cfg.DatabaseDSN)
		if err != nil {
			return err
		}

		srv.store = store
		return nil
	}

	store := storage.NewMemStorage(srv.cfg.StoreInternal == 0, srv.cfg.FileStoragePath)
	if srv.cfg.RestoreData {
		err := store.LoadStorageFromFile()
		if err != nil {
			slog.Warn("restore data at start server", slog.String("error", err.Error()))
		}
	}
	srv.store = store

	return nil
}

func (srv *Server) migrate() error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", srv.cfg.DatabaseDSN)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&migrations.Content)
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	return db.Close()
}

func (srv *Server) startHTTPServer() error {
	slog.Info(fmt.Sprintf("Running server on %s", srv.cfg.ServerAddr))
	err := srv.httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (srv *Server) shutdown(ctx context.Context) error {
	err := srv.httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown http server")
		return err
	}

	err = srv.store.DumpStorageToFile()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = srv.store.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Server shutdown complete")

	return err
}

func (srv *Server) Start(ctx context.Context) error {

	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return srv.startHTTPServer()
	})

	grp.Go(func() error {
		<-ctx.Done()

		return srv.shutdown(ctx)
	})

	if srv.cfg.StoreInternal > 0 {
		grp.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(srv.cfg.StoreInternal):
					err := utils.RetryFunc(ctx, func() error {
						return srv.store.DumpStorageToFile()
					})
					if err != nil {
						slog.Error(fmt.Sprintf("error autosave: %s", err))
					}
				}
			}
		})
	}

	return grp.Wait()
}
