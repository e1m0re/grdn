package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage"
)

type metricService struct {
	store storage.StoreManager
}

func NewMetricsService(store storage.StoreManager) MetricsService {
	return &metricService{
		store: store,
	}
}

func (ms *metricService) PingDB(ctx context.Context) error {
	return ms.store.Ping(ctx)
}
func (ms *metricService) GetMetricsList(ctx context.Context) ([]string, error) {
	return ms.store.GetMetricsList(ctx)
}
func (ms *metricService) GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error) {
	return ms.store.GetMetric(ctx, mType, mName)
}
func (ms *metricService) UpdateMetric(ctx context.Context, metric models.Metric) error {
	return ms.store.UpdateMetric(ctx, metric)
}
func (ms *metricService) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	return ms.store.UpdateMetrics(ctx, metrics)
}
