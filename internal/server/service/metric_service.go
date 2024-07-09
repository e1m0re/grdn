package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage"
)

type metricService struct {
	store storage.Storage
}

// NewMetricService is metricService constructor.
func NewMetricService(store storage.Storage) MetricService {
	return &metricService{
		store: store,
	}
}

// GetMetric returns an object Metric.
func (ms *metricService) GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (metric *models.Metric, err error) {
	return ms.store.GetMetric(ctx, mType, mName)
}

// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
func (ms *metricService) GetMetricsList(ctx context.Context) ([]string, error) {
	return ms.store.GetMetricsList(ctx)
}

// UpdateMetric performs updates to the value of the specified result in the store.
func (ms *metricService) UpdateMetric(ctx context.Context, metric models.Metric) error {
	return ms.store.UpdateMetric(ctx, metric)
}

// UpdateMetrics performs batch updates of result values in the store.
func (ms *metricService) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	return ms.store.UpdateMetrics(ctx, metrics)
}
