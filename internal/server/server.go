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
}

// NewServer is Server constructor.
func NewServer(cfg *config.Config) *Server {
	srv := &Server{
		cfg: cfg,
	}

	services := service.NewServices()
	handler := appHandler.NewHandler(services)
	srv.httpServer = &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: handler.NewRouter(cfg.Key),
	}

	return srv
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

	err = store.Get().Save(ctx)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = store.Get().Close()
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
