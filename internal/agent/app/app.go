package app

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"

	"github.com/e1m0re/grdn/internal/agent/apiclient"
	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/agent/monitor"
)

type App struct {
	apiClient *apiclient.APIClient
	cfg       *config.Config
	monitor   *monitor.MetricsMonitor
}

func NewApp(cfg *config.Config) *App {
	return &App{
		apiClient: apiclient.NewAPI("http://" + cfg.ServerAddr),
		cfg:       cfg,
		monitor:   monitor.NewMetricsMonitor(),
	}
}

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
			case <-time.After(app.cfg.ReportInterval):
				app.sendDataToServer(ctx)
			}
		}
	})

	return grp.Wait()
}

func (app *App) sendDataToServer(ctx context.Context) {
	metrics := app.monitor.GetMetricsList()

	content, err := json.Marshal(metrics)
	if err != nil {
		slog.Error("Error marshalling metrics data",
			slog.String("error", err.Error()),
		)
		return
	}

	err = retryFunc(ctx, func() error {
		return app.apiClient.SendMetricsData(&content)
	})
	if err != nil {
		slog.Error("send metrics data failed",
			slog.String("error", err.Error()),
		)
	}
}
