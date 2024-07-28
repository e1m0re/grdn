// Package app implements initialize and start clients application.
package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/service"
	"github.com/e1m0re/grdn/internal/service/apiclient"
	"github.com/e1m0re/grdn/internal/service/encryption"
	"github.com/e1m0re/grdn/internal/service/monitor"
	"github.com/e1m0re/grdn/internal/utils"
)

type contentType = []byte

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=App
type App interface {
	// Start runs client application.
	Start(ctx context.Context) error
}

type app struct {
	apiClient apiclient.APIClient
	cfg       *config.Config
	monitor   monitor.Monitor
	encryptor encryption.Encryptor
}

// Start runs client application.
func (app *app) Start(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return app.updateDataWorker(ctx)
	})

	grp.Go(func() error {
		return app.updateGOPSDataWorker(ctx)
	})

	tasksQueue := make(chan contentType, 10)
	defer close(tasksQueue)

	grp.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				app.sendDataToServer(tasksQueue)
				return nil
			case <-time.After(app.cfg.ReportInterval):
				app.sendDataToServer(tasksQueue)
			}
		}
	})

	for i := 1; i <= app.cfg.RateLimit; i++ {
		grp.Go(func() error {
			return app.sendDataToServerWorker(ctx, tasksQueue)
		})
	}

	return grp.Wait()
}

func (app *app) updateDataWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(app.cfg.PollInterval):
			app.monitor.UpdateData()
		}
	}
}

func (app *app) updateGOPSDataWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(app.cfg.PollInterval):
			err := app.monitor.UpdateGOPS(ctx)
			if err != nil {
				slog.Error("error update GOPS data", slog.String("error", err.Error()))
			}
		}
	}
}

func (app *app) sendDataToServer(outChan chan<- contentType) {
	metrics := app.monitor.GetMetricsList()

	content, err := json.Marshal(metrics)
	if err != nil {
		slog.Error("Error marshalling metrics data",
			slog.String("error", err.Error()),
		)
		return
	}

	if app.encryptor != nil {
		content, err = app.encryptor.Encrypt(content)
		if err != nil {
			slog.Error("encryption error", slog.String("error", err.Error()))
		}
	}

	outChan <- content
}

func (app *app) sendDataToServerWorker(ctx context.Context, tasksQueue <-chan contentType) error {
	for {
		select {
		case <-ctx.Done():
		case c, ok := <-tasksQueue:
			if !ok {
				return nil
			}
			err := utils.RetryFunc(ctx, func() error {
				return app.apiClient.SendMetricsData(&c)
			})

			if err != nil {
				slog.Error("send metrics data failed", slog.String("error", err.Error()))
			}
		}
	}
}

// NewApp is app constructor.
func NewApp(cfg *config.Config, services *service.AgentServices) App {
	return &app{
		apiClient: services.APIClient,
		cfg:       cfg,
		monitor:   services.Monitor,
		encryptor: services.Encryptor,
	}
}
