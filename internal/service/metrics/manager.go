package metrics

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage"
	"github.com/e1m0re/grdn/internal/storage/store"
)

// Manager is the interface that contains all operations for metrics.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Manager
type Manager interface {
	// GetAllMetrics returns result of all metrics.
	GetAllMetrics(ctx context.Context) (*models.MetricsList, error)

	// GetMetric returns an object Metric. Returns nil,nil if metric not found.
	GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (*models.Metric, error)

	// UpdateMetric performs updates to the value of the specified result in the store.
	UpdateMetric(ctx context.Context, metric models.Metric) error

	// UpdateMetrics performs batch updates of result values in the store.
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
}

type metricsManager struct {
	store store.Store
}

// NewMetricsManager returns new instance of metrics manager.
func NewMetricsManager(s store.Store) Manager {
	return &metricsManager{
		store: s,
	}
}

// GetAllMetrics returns result of all metrics.
func (mm *metricsManager) GetAllMetrics(ctx context.Context) (*models.MetricsList, error) {
	return mm.store.GetAllMetrics(ctx)
}

// GetMetric returns an object Metric. Returns nil,nil if metric not found.
func (mm *metricsManager) GetMetric(ctx context.Context, mType models.MetricType, mName models.MetricName) (*models.Metric, error) {
	return mm.store.GetMetric(ctx, mType, mName)
}

func (mm *metricsManager) processUpdateMetric(ctx context.Context, metric models.Metric) (*models.Metric, error) {
	cm, err := mm.store.GetMetric(ctx, metric.MType, metric.ID)
	if err != nil {
		return nil, err
	}
	if cm == nil {
		cm = &models.Metric{
			Value: nil,
			Delta: nil,
			MType: metric.MType,
			ID:    metric.ID,
		}
	}

	switch cm.MType {
	case models.GaugeType:
		cm.Value = metric.Value
	case models.CounterType:
		if cm.Delta != nil {
			newDelta := *cm.Delta + *metric.Delta
			cm.Delta = &newDelta
		} else {
			cm.Delta = metric.Delta
		}
	default:
		return nil, storage.ErrUnknownMetricType
	}

	return cm, nil
}

// UpdateMetric performs updates to the value of the specified result in the store.
func (mm *metricsManager) UpdateMetric(ctx context.Context, metric models.Metric) error {
	cm, err := mm.processUpdateMetric(ctx, metric)
	if err != nil {
		return err
	}

	return mm.store.UpdateMetrics(ctx, models.MetricsList{cm})
}

// UpdateMetrics performs batch updates of result values in the store.
func (mm *metricsManager) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	if len(metrics) == 0 {
		return nil
	}

	for i, metric := range metrics {
		cm, err := mm.processUpdateMetric(ctx, *metric)
		if err != nil {
			return err
		}

		metrics[i] = cm
	}

	return mm.store.UpdateMetrics(ctx, metrics)
}
