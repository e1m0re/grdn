package storage

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
)

// Store is an interface that defines methods for working with storage.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Store
type Store interface {
	// Close closes the connection to the storage.
	Close() error

	// DumpStorageToFile saves data to a file.
	DumpStorageToFile() error

	// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
	GetMetricsList(ctx context.Context) ([]string, error)

	// GetMetric returns an object Metric.
	GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error)

	// LoadStorageFromFile loads data from a file.
	LoadStorageFromFile() error

	// Ping checks the connection to the storage.
	Ping(ctx context.Context) error

	// UpdateMetric performs updates to the value of the specified metric in the store.
	UpdateMetric(ctx context.Context, metric models.Metric) error

	// UpdateMetrics performs batch updates of metric values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}
