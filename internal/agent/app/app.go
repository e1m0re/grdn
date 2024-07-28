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

type content = []byte

type App struct {
	apiClient *apiclient.APIClient
	cfg       *config.Config
	monitor   monitor.Monitor
	encryptor encryption.Encryptor
}

// NewApp is App constructor.
func NewApp(cfg *config.Config, services *service.AgentServices) *App {
	return &App{
		apiClient: services.APIClient,
		cfg:       cfg,
		monitor:   services.Monitor,
		encryptor: services.Encryptor,
	}
}

// Start runs client application.
func (app *App) Start(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return app.updateDataWorker(ctx)
	})

	grp.Go(func() error {
		return app.updateGOPSData(ctx)
	})

	tasksQueue := make(chan content, 10)
	defer close(tasksQueue)

	grp.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(app.cfg.ReportInterval):
				app.sendDataToServer(ctx, tasksQueue)
			}
		}
	})

	grp.Go(func() error {
		<-ctx.Done()
		app.sendDataToServer(ctx, tasksQueue)

		return nil
	})

	for i := 1; i <= app.cfg.RateLimit; i++ {
		grp.Go(func() error {
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
						slog.Error("send metrics data failed",
							slog.String("error", err.Error()),
						)
					}
				}
			}
		})
	}

	return grp.Wait()
}

func (app *App) updateDataWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(app.cfg.PollInterval):
			app.monitor.UpdateData()
		}
	}
}

func (app *App) updateGOPSData(ctx context.Context) error {
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

func (app *App) sendDataToServer(ctx context.Context, outChan chan<- content) {
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
