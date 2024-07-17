// Package app implements initialize and start clients application.
package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/e1m0re/grdn/internal/agent/apiclient"
	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/agent/monitor"
	"github.com/e1m0re/grdn/internal/encryption"
	"github.com/e1m0re/grdn/internal/utils"
)

type content = []byte

type App struct {
	apiClient *apiclient.APIClient
	cfg       *config.Config
	monitor   *monitor.MetricsMonitor
	encryptor encryption.Encryptor
}

// NewApp is App constructor.
func NewApp(cfg *config.Config) *App {
	var encr encryption.Encryptor
	var err error
	if len(cfg.PublicKeyFile) > 0 {
		encr, err = encryption.NewEncryptor(cfg.PublicKeyFile)
		if err != nil {
			panic(err)
		}
	}

	return &App{
		apiClient: apiclient.NewAPIClient("http://"+cfg.ServerAddr, []byte(cfg.Key)),
		cfg:       cfg,
		monitor:   monitor.NewMetricsMonitor(),
		encryptor: encr,
	}
}

// Start runs client application.
func (app *App) Start(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(app.cfg.PollInterval):
				app.monitor.UpdateData()
			}
		}
	})

	grp.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(app.cfg.PollInterval):
				app.monitor.UpdateGOPS(ctx)
			}
		}
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
