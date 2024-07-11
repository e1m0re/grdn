package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage"
	"github.com/e1m0re/grdn/internal/server/storage/store/memory"
	"github.com/e1m0re/grdn/internal/server/storage/store/sql"
	"github.com/e1m0re/grdn/internal/utils"
)

// Store is the interface that each store should implement
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Store
type Store interface {

	// Close closes the connection to the storage.
	Close() error

	// GetAllMetrics returns the list of all metrics.
	GetAllMetrics(ctx context.Context) (*models.MetricsList, error)

	// GetMetric returns an object Metric.
	GetMetric(ctx context.Context, mType models.MetricType, mName string) (*models.Metric, error)

	// Ping checks the connection to the storage.
	Ping(ctx context.Context) error

	// Restore loads data from a file.
	Restore(ctx context.Context) error

	// Save saves data to a file.
	Save(ctx context.Context) error

	// UpdateMetrics performs batch updates of result values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

var (
	store Store

	initialized = false

	ctx        context.Context
	cancelFunc context.CancelFunc
)

// Get returns link to store.
func Get() Store {
	if !initialized {
		err := Initialize(nil)
		if err != nil {
			panic("failed auto initialization of store " + err.Error())
		}
	}

	return store
}

// Initialize instantiates the storage provider based on the Config provider
func Initialize(cfg *storage.Config) error {
	initialized = true
	var err error

	if cancelFunc != nil {
		// Stop the active autoSave task, if there's already one
		cancelFunc()
	}
	ctx, cancelFunc = context.WithCancel(context.Background())

	switch cfg.Type {
	case storage.TypePostgres:
		store, err = sql.NewStore(cfg.Path)
	case storage.TypeMemory:
		fallthrough
	default:
		store, _ = memory.NewStore(context.Background(), cfg.Path, cfg.SyncMode)
		autoSave(ctx, store, cfg.Interval)
	}

	return err
}

// autoSave automatically calls the Save function of the provider at every interval
func autoSave(ctx context.Context, store Store, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("[store.autoSave] Stopping active job")
			return
		case <-time.After(interval):
			slog.Info("[store.autoSave] Saving")
			err := utils.RetryFunc(ctx, func() error {
				return store.Save(ctx)
			})
			if err != nil {
				slog.Info("[store.autoSave] Save failed:", "error", err.Error())
			}
		}
	}
}
