package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/server/config"
	gzipMiddleware "github.com/e1m0re/grdn/internal/server/middleware/gzip"
	loggerMiddleware "github.com/e1m0re/grdn/internal/server/middleware/logger"
	"github.com/e1m0re/grdn/internal/storage"
)

type Server struct {
	cfg        *config.Config
	router     *chi.Mux
	httpServer *http.Server
	store      storage.Interface
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		slog.Error(fmt.Sprintf("error init database connection: %s", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		slog.Error(fmt.Sprintf("error init database connection: %s", err))
	}

	store := storage.NewMemStorage(cfg.StoreInternal == 0, cfg.FileStoragePath)
	if cfg.RestoreData {
		err := store.LoadStorageFromFile()
		if err != nil {
			slog.Error(fmt.Sprintf("error restore data: %s", err))
		}
	}

	srv := &Server{
		cfg:    cfg,
		router: chi.NewRouter(),
		store:  store,
	}

	srv.initRoutes()

	return srv, nil
}

func (srv *Server) initRoutes() {
	srv.router.Use(loggerMiddleware.Middleware)
	srv.router.Use(gzipMiddleware.Middleware)
	srv.router.Use(middleware.Compress(5, "text/html", "application/json"))

	srv.router.Route("/", func(r chi.Router) {
		r.Get("/", srv.getMainPage)
		r.Get("/ping", srv.checkDBConnection)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", srv.getMetricValueV2)
			r.Get("/{mType}/{mName}", srv.getMetricValue)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", srv.updateMetrics)
			r.Post("/{mType}/{mName}/{mValue}", srv.updateMetric)
		})
	})
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
	}

	err = srv.store.DumpStorageToFile()
	if err != nil {
		slog.Error(err.Error())
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
					err := srv.store.DumpStorageToFile()
					if err != nil {
						slog.Error(fmt.Sprintf("error autosave: %s", err))
					}
				}
			}
		})
	}

	return grp.Wait()
}
