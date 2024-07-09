package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=MetricService
type MetricService interface {
	// GetMetric returns an object Metric.
	GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (metric *models.Metric, err error)

	// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
	GetMetricsList(ctx context.Context) ([]string, error)

	// UpdateMetric performs updates to the value of the specified result in the store.
	UpdateMetric(ctx context.Context, metric models.Metric) error

	// UpdateMetrics performs batch updates of result values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=StorageService
type StorageService interface {
	// Close closes the connection to the storage.
	Close() error

	// DumpStorageToFile saves data to a file.
	DumpStorageToFile() error

	// LoadStorageFromFile loads data from a file.
	LoadStorageFromFile() error

	// PingDB checks the connection to the storage.
	PingDB(ctx context.Context) error
}

type Services struct {
	MetricService
	StorageService
}

// NewServices is Services constructor.
func NewServices(store storage.Storage) *Services {
	return &Services{
		MetricService:  NewMetricService(store),
		StorageService: NewStorageService(store),
	}
}
