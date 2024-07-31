package store

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage"
	"github.com/e1m0re/grdn/internal/storage/store/memory"
	"github.com/e1m0re/grdn/internal/storage/store/sql"
	"github.com/e1m0re/grdn/internal/utils"
)

// Store is the interface that each store should implement
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Store
type Store interface {

	// Clear removes all data in storage.
	Clear(ctx context.Context) error

	// Close closes the connection to the storage.
	Close() error

	// GetAllMetrics returns the list of all metrics.
	GetAllMetrics(ctx context.Context) (*models.MetricsList, error)

	// GetMetric returns an object Metric. Returns nil,nil if metric not found.
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

// NewStore instantiates the storage provider based on the Config provider
func NewStore(ctx context.Context, cfg *storage.Config) (Store, error) {
	if cfg == nil {
		// This only happens in tests
		log.Println("[store.Initialize] nil storage config passed as parameter. This should only happen in tests. Defaulting to an empty config.")
		cfg = &storage.Config{}
	}

	var (
		store Store
		err   error
	)
	switch cfg.Type {
	case storage.TypePostgres:
		store, err = sql.NewStore("pgx", cfg.Path)
	case storage.TypeMemory:
		fallthrough
	default:
		store, _ = memory.NewStore(ctx, cfg.Path, cfg.SyncMode)
		go autoSave(ctx, store, cfg.Interval)
	}

	return store, err
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
