package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"

	"golang.org/x/sync/errgroup"

	appHandler "github.com/e1m0re/grdn/internal/api"
	"github.com/e1m0re/grdn/internal/listeners"
	"github.com/e1m0re/grdn/internal/listeners/http/config"
	"github.com/e1m0re/grdn/internal/service"
	"github.com/e1m0re/grdn/internal/storage/store"
)

type listener struct {
	cfg        *config.Config
	httpServer *http.Server
	services   *service.ServerServices
}

// Run starts HTTP listener.
func (l *listener) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return l.startHTTPServer()
	})

	grp.Go(func() error {
		<-ctx.Done()

		return l.Shutdown(ctx)
	})

	return grp.Wait()
}

// Shutdown stops HTTP listener.
func (l *listener) Shutdown(ctx context.Context) error {
	err := l.httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown http server")
		return err
	}

	err = l.services.StorageService.Save(ctx)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = l.services.StorageService.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("listener shutdown complete")

	return err
}

func (l *listener) startHTTPServer() error {
	slog.Info(fmt.Sprintf("Running server on %s", l.cfg.ServerAddr))
	err := l.httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

var _ listeners.Listener = (*listener)(nil)

// NewHTTPListener initiates new instance of HTTP listener.
func NewHTTPListener(cfg *config.Config, s store.Store) listeners.Listener {
	services := service.NewServerServices(s)

	trustedSubnet, _ := netip.ParsePrefix(cfg.TrustedSubnet)
	handlerConfig := appHandler.Config{
		SignKey:        cfg.Key,
		PrivateKeyFile: cfg.PrivateKeyFile,
		TrustedSubnet:  &trustedSubnet,
	}
	handler := appHandler.NewHandler(services, handlerConfig)

	return &listener{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.ServerAddr,
			Handler: handler.NewRouter(),
		},
		services: services,
	}
}
