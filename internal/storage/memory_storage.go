package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e1m0re/grdn/internal/models"
	"log/slog"
	"os"
	"strconv"
)

type MemStorage struct {
	Gauges   map[GaugeName]GaugeDateType
	Counters map[CounterName]CounterDateType
	syncMode bool
	filePath string
}

func NewMemStorage(syncMode bool, filePath string) *MemStorage {
	return &MemStorage{
		Gauges:   make(map[GaugeName]GaugeDateType),
		Counters: make(map[CounterName]CounterDateType),
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

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
func (s *MemStorage) GetAllMetrics(ctx context.Context) ([]string, error) {
	var result []string
	for key, value := range s.Gauges {
		result = append(result, fmt.Sprintf("%s: %s", key, strconv.FormatFloat(value, 'f', -1, 64)))
	}

	for key, value := range s.Counters {
		result = append(result, fmt.Sprintf("%s: %v", key, value))
	}

	return result, nil
}
func (s *MemStorage) GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error) {
	switch mType {
	case models.GaugeType:
		value, ok := s.Gauges[mName]

		if !ok {
			err = errors.New("unknown metric")
			return nil, err
		}

		metric = &models.Metric{
			ID:    mName,
			MType: models.GaugeType,
			Delta: nil,
			Value: &value,
		}
	case models.CounterType:
		delta, ok := s.Counters[mName]

		if !ok {
			err = errors.New("unknown metric")
			return nil, err
		}

		metric = &models.Metric{
			ID:    mName,
			MType: models.CounterType,
			Delta: &delta,
			Value: nil,
		}
	default:
		err = errors.New("unknown metric")
	}

	return
}
func (s *MemStorage) LoadStorageFromFile() error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	tmpData := &MemStorage{}
	err = json.Unmarshal(file, &tmpData)

	if err == nil {
		s.Gauges = tmpData.Gauges
		s.Counters = tmpData.Counters
	}

	return err
}
func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}
func (s *MemStorage) UpdateMetric(ctx context.Context, metric models.Metric) error {
	switch metric.MType {
	case models.GaugeType:
		s.Gauges[metric.ID] = *metric.Value
	case models.CounterType:
		s.Counters[metric.ID] = *metric.Delta
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
