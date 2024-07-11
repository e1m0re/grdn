package metrics

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage/store"
)

// Manager is the interface that contains all operations for metrics.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Manager
type Manager interface {
	// GetAll returns list of all metrics.
	GetAll(ctx context.Context) (*models.MetricsList, error)

	// GetMetric returns an object Metric.
	GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (*models.Metric, error)

	// GetSimpleMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
	GetSimpleMetricsList(ctx context.Context) ([]string, error)

	// UpdateMetric performs updates to the value of the specified result in the store.
	UpdateMetric(ctx context.Context, metric models.Metric) error

	// UpdateMetrics performs batch updates of result values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

type metricsManager struct{}

// NewMetricsManager returns new instance of metrics manager.
func NewMetricsManager() Manager {
	return &metricsManager{}
}

// GetAll returns list of all metrics.
func (mm *metricsManager) GetAll(ctx context.Context) (*models.MetricsList, error) {
	return store.Get().GetAllMetrics(ctx)
}

// GetMetric returns an object Metric.
func (mm *metricsManager) GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (*models.Metric, error) {
	return store.Get().GetMetric(ctx, mType, mName)
}

// GetSimpleMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
func (mm *metricsManager) GetSimpleMetricsList(ctx context.Context) ([]string, error) {
	metricsList, err := mm.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(*metricsList))
	for i, metric := range *metricsList {
		result[i] = metric.String()
	}
	return result, nil
}

// UpdateMetric performs updates to the value of the specified result in the store.
func (mm *metricsManager) UpdateMetric(ctx context.Context, metric models.Metric) error {
	//switch metric.MType {
	//case models.GaugeType:
	//	s.gauges[metric.ID] = *metric.Value
	//case models.CounterType:
	//	s.counters[metric.ID] += *metric.Delta
	//default:
	//	return errors.New("unknown metric type")
	//}
	//
	//return store.Get().UpdateMetric(ctx, metric)

	return nil
}

// UpdateMetrics performs batch updates of result values in the store.
func (mm *metricsManager) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	return store.Get().UpdateMetrics(ctx, metrics)
}
