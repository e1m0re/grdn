package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=MetricsService
type MetricsService interface {
	// PingDB checks the connection to the storage.
	PingDB(ctx context.Context) error

	// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
	GetMetricsList(ctx context.Context) ([]string, error)

	// GetMetric returns an object Metric.
	GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error)

	// UpdateMetric performs updates to the value of the specified metric in the store.
	UpdateMetric(ctx context.Context, metric models.Metric) error

	// UpdateMetrics performs batch updates of metric values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

type Services struct {
	MetricsService
}

func NewServices(store storage.Store) *Services {
	return &Services{
		MetricsService: NewMetricsService(store),
	}
}
