package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/e1m0re/grdn/internal/models"
)

type metrics struct {
	Gauges   map[models.GaugeName]models.GaugeDateType     `json:"gauges"`
	Counters map[models.CounterName]models.CounterDateType `json:"counters"`
}

type MemStorage struct {
	gauges   map[models.GaugeName]models.GaugeDateType
	counters map[models.CounterName]models.CounterDateType
	filePath string
	syncMode bool
	mx       sync.RWMutex
}

// NewMemStorage is MemStorage constructor.
func NewMemStorage(syncMode bool, filePath string) *MemStorage {
	return &MemStorage{
		gauges:   make(map[models.GaugeName]models.GaugeDateType),
		counters: make(map[models.CounterName]models.CounterDateType),
		syncMode: syncMode,
		filePath: filePath,
	}
}

// Close closes the connection to the storage.
func (s *MemStorage) Close() error {
	return nil
}

// DumpStorageToFile saves data to a file.
func (s *MemStorage) DumpStorageToFile() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	metrics := &metrics{
		Gauges:   s.gauges,
		Counters: s.counters,
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

// GetMetric returns an object Metric.
func (s *MemStorage) GetMetric(ctx context.Context, mType models.MetricType, mName string) (metric *models.Metric, err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	switch mType {
	case models.GaugeType:
		value, ok := s.gauges[mName]

		if !ok {
			return nil, ErrUnknownMetric
		}

		metric = &models.Metric{
			ID:    mName,
			MType: models.GaugeType,
			Delta: nil,
			Value: &value,
		}
	case models.CounterType:
		delta, ok := s.counters[mName]

		if !ok {
			return nil, ErrUnknownMetric
		}

		metric = &models.Metric{
			ID:    mName,
			MType: models.CounterType,
			Delta: &delta,
			Value: nil,
		}
	default:

		return nil, ErrUnknownMetric
	}

	return
}

// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
func (s *MemStorage) GetMetricsList(ctx context.Context) ([]string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	result := make([]string, len(s.gauges)+len(s.counters))
	i := 0
	for key, value := range s.gauges {
		result[i] = fmt.Sprintf("%s: %s", key, strconv.FormatFloat(value, 'f', -1, 64))
		i++
	}

	for key, value := range s.counters {
		result[i] = fmt.Sprintf("%s: %v", key, value)
		i++
	}

	return result, nil
}

// LoadStorageFromFile loads data from a file.
func (s *MemStorage) LoadStorageFromFile() error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	m := metrics{
		Gauges:   make(map[models.GaugeName]models.GaugeDateType),
		Counters: make(map[models.CounterName]models.CounterDateType),
	}
	err = json.Unmarshal(file, &m)

	if err == nil {
		s.gauges = m.Gauges
		s.counters = m.Counters
	}

	return err
}

// Ping checks the connection to the storage.
func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}

// UpdateMetric performs updates to the value of the specified metric in the store.
func (s *MemStorage) UpdateMetric(ctx context.Context, metric models.Metric) error {
	switch metric.MType {
	case models.GaugeType:
		s.gauges[metric.ID] = *metric.Value
	case models.CounterType:
		s.counters[metric.ID] += *metric.Delta
	default:
		return errors.New("unknown metric type")
	}

	if s.syncMode {
		err := s.DumpStorageToFile()
		if err != nil {
			slog.Error(fmt.Sprintf("error on save data to HDD: %s", err))
		}
	}

	return nil
}

// UpdateMetrics performs batch updates of metric values in the store.
func (s *MemStorage) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	for _, metric := range metrics {
		err := s.UpdateMetric(ctx, *metric)
		if err != nil {
			return err
		}
	}

	return nil
}
