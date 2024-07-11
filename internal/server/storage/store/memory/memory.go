package memory

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/utils"
)

// Store that leverages RAM.
type Store struct {
	metrics map[string]models.Metric

	filePath string
	syncMode bool
	sync.RWMutex
}

// NewStore creates a new in-memory store.
func NewStore(ctx context.Context, filePath string, syncMode bool) (*Store, error) {
	store := &Store{
		metrics:  make(map[string]models.Metric),
		syncMode: syncMode,
		filePath: filePath,
	}

	var err error
	if len(filePath) > 0 {
		err = store.Restore(ctx)
	}

	return store, err
}

func (s *Store) genMetricKey(m models.MetricName, t models.MetricType) string {
	return utils.GetMD5Hash(t + m)
}

// Clear removes all data in storage.
func (s *Store) Clear(ctx context.Context) error {
	s.metrics = make(map[string]models.Metric)
	return nil
}

// Close closes the connection to the storage.
func (s *Store) Close() error {
	return nil
}

// GetAllMetrics returns the list of all metrics.
func (s *Store) GetAllMetrics(ctx context.Context) (*models.MetricsList, error) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	result := make(models.MetricsList, len(s.metrics))
	i := 0
	for _, metric := range s.metrics {
		result[i] = &models.Metric{
			Value: metric.Value,
			Delta: metric.Delta,
			MType: metric.MType,
			ID:    metric.ID,
		}

		i++
	}

	return &result, nil
}

// GetMetric returns an object Metric.
func (s *Store) GetMetric(ctx context.Context, mType models.MetricType, mName string) (*models.Metric, error) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	key := s.genMetricKey(mName, mType)
	metric, ok := s.metrics[key]
	if !ok {
		return nil, nil
	}

	return &metric, nil
}

// Ping checks the connection to the storage.
func (s *Store) Ping(ctx context.Context) error {
	return nil
}

// Restore loads data from a file.
func (s *Store) Restore(ctx context.Context) error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var metrics models.MetricsList
	err = json.Unmarshal(file, &metrics)
	if err != nil {
		return err
	}

	s.metrics = make(map[string]models.Metric, len(metrics))

	return s.UpdateMetrics(ctx, metrics)
}

// Save saves data to a file.
func (s *Store) Save(ctx context.Context) error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	metrics, _ := s.GetAllMetrics(ctx)
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

// UpdateMetrics performs batch updates of result values in the store.
func (s *Store) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	for _, metric := range metrics {
		key := s.genMetricKey(metric.ID, metric.MType)
		s.metrics[key] = *metric
	}

	if s.syncMode {
		err := s.Save(ctx)

		return err
	}

	return nil
}
