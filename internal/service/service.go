package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=MetricsService
type MetricsService interface {
	PingDB(ctx context.Context) error
	GetMetricsList(ctx context.Context) ([]string, error)
	GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error)
	UpdateMetric(ctx context.Context, metric models.Metric) error
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

type Services struct {
	MetricsService
}

func NewServices(store storage.StoreManager) *Services {
	return &Services{
		MetricsService: NewMetricsService(store),
	}
}
