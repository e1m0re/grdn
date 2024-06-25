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
	Gauges   map[GaugeName]GaugeDateType     `json:"gauges"`
	Counters map[CounterName]CounterDateType `json:"counters"`
}

type MemStorage struct {
	mx       sync.RWMutex
	gauges   map[GaugeName]GaugeDateType
	counters map[CounterName]CounterDateType
	syncMode bool
	filePath string
}

func NewMemStorage(syncMode bool, filePath string) *MemStorage {
	return &MemStorage{
		gauges:   make(map[GaugeName]GaugeDateType),
		counters: make(map[CounterName]CounterDateType),
		syncMode: syncMode,
		filePath: filePath,
	}
}

func (s *MemStorage) Close() error {
	return nil
}
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
func (s *MemStorage) GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error) {
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
func (s *MemStorage) LoadStorageFromFile() error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	metrics := metrics{
		Gauges:   make(map[GaugeName]GaugeDateType),
		Counters: make(map[CounterName]CounterDateType),
	}
	err = json.Unmarshal(file, &metrics)

	if err == nil {
		s.gauges = metrics.Gauges
		s.counters = metrics.Counters
	}

	return err
}
func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}
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
func (s *MemStorage) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	for _, metric := range metrics {
		err := s.UpdateMetric(ctx, *metric)
		if err != nil {
			return err
		}
	}

	return nil
}
