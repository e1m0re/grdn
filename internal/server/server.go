package server

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"net/netip"

	appHandler "github.com/e1m0re/grdn/internal/api"
	"github.com/e1m0re/grdn/internal/server/config"
	"github.com/e1m0re/grdn/internal/service"
	"github.com/e1m0re/grdn/internal/storage/store"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Server
type Server interface {
	// Start runs server.
	Start(ctx context.Context) error
}

type srv struct {
	cfg        *config.Config
	httpServer *http.Server
	services   *service.ServerServices
}

// Start runs server.
func (srv *srv) Start(ctx context.Context) error {
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

func (srv *srv) startHTTPServer() error {
	slog.Info(fmt.Sprintf("Running server on %s", srv.cfg.ServerAddr))
	err := srv.httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (srv *srv) shutdown(ctx context.Context) error {
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

	slog.Info("srv shutdown complete")

	return err
}

// NewServer is srv constructor.
func NewServer(cfg *config.Config, s store.Store) Server {
	services := service.NewServerServices(s)

	trustedSubnet, _ := netip.ParsePrefix(cfg.TrustedSubnet)
	handlerConfig := appHandler.Config{
		SignKey:        cfg.Key,
		PrivateKeyFile: cfg.PrivateKeyFile,
		TrustedSubnet:  &trustedSubnet,
	}
	handler := appHandler.NewHandler(services, handlerConfig)

	return &srv{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.ServerAddr,
			Handler: handler.NewRouter(),
		},
		services: services,
	}
}
