package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/db/migrations"
	"github.com/e1m0re/grdn/internal/server/config"
	grdnMiddleware "github.com/e1m0re/grdn/internal/server/middleware"
	"github.com/e1m0re/grdn/internal/storage"
)

type Server struct {
	cfg        *config.Config
	router     *chi.Mux
	httpServer *http.Server
	store      storage.Interface
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	srv := &Server{
		cfg:    cfg,
		router: chi.NewRouter(),
	}

	srv.initRoutes()

	err := srv.initStore(ctx)

	return srv, err
}

func (srv *Server) initStore(ctx context.Context) error {
	if srv.cfg.DatabaseDSN != "" {
		err := srv.migrate(ctx)
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

func (srv *Server) initRoutes() {
	srv.router.Use(grdnMiddleware.LoggingMiddleware)
	srv.router.Use(grdnMiddleware.GZipMiddleware)
	srv.router.Use(middleware.Compress(5, "text/html", "application/json"))

	srv.router.Route("/", func(r chi.Router) {
		r.Get("/", srv.getMainPage)
		r.Get("/ping", srv.checkDBConnection)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", srv.getMetricValueV2)
			r.Get("/{mType}/{mName}", srv.getMetricValue)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", srv.updateMetricV2)
			r.Post("/{mType}/{mName}/{mValue}", srv.updateMetricV1)
		})
		r.Route("/updates", func(r chi.Router) {
			r.Post("/", srv.updateMetrics)
		})
	})
}

func (srv *Server) migrate(ctx context.Context) error {
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

func (srv *Server) startHTTPServer(ctx context.Context) error {
	srv.httpServer = &http.Server{
		Addr:    srv.cfg.ServerAddr,
		Handler: srv.router,
	}

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
		return srv.startHTTPServer(ctx)
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
					err := retryFunc(ctx, func() error {
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
