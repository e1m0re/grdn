package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/server/config"
	appHandler "github.com/e1m0re/grdn/internal/server/handler"
	"github.com/e1m0re/grdn/internal/server/service"
	"github.com/e1m0re/grdn/internal/server/storage/store"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	services   *service.Services
}

// NewServer is Server constructor.
func NewServer(cfg *config.Config, s store.Store) *Server {
	services := service.NewServices(s)
	handler := appHandler.NewHandler(services)

	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.ServerAddr,
			Handler: handler.NewRouter(cfg.Key),
		},
		services: services,
	}
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

	err = srv.services.StorageService.Save(ctx)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = srv.services.StorageService.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Server shutdown complete")

	return err
}

// Start runs server.
func (srv *Server) Start(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return srv.startHTTPServer()
	})

	grp.Go(func() error {
		<-ctx.Done()

		return srv.shutdown(ctx)
	})

	return grp.Wait()
}
