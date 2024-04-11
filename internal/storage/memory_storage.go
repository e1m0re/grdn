package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/e1m0re/grdn/internal/models"
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

func (s *MemStorage) UpdateGaugeMetric(name GaugeName, value GaugeDateType) {
	s.Gauges[name] = value
}

func (s *MemStorage) UpdateCounterMetric(name CounterName, value CounterDateType) {
	s.Counters[name] += value
}

func (s *MemStorage) UpdateMetricValue(mType models.MetricsType, mName string, mValue string) error {
	switch mType {
	case models.GaugeType:
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return err
		}

		s.UpdateGaugeMetric(mName, value)
	case models.CounterType:
		value, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return err
		}

		s.UpdateCounterMetric(mName, value)
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

func (s *MemStorage) UpdateMetricValueV2(data models.Metrics) error {
	if len(data.ID) == 0 {
		return errors.New("invalid metrics name")
	}
	switch data.MType {
	case models.GaugeType:
		if data.Value == nil {
			return errors.New("invalid value")
		}
		s.UpdateGaugeMetric(data.ID, *data.Value)
	case models.CounterType:
		if data.Delta == nil {
			return errors.New("invalid value")
		}
		s.UpdateCounterMetric(data.ID, *data.Delta)
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

func (s *MemStorage) GetMetric(mType models.MetricsType, mName string) (metric *models.Metrics, err error) {
	switch mType {
	case models.GaugeType:
		value, ok := s.Gauges[mName]

		if !ok {
			err = errors.New("unknown metric")
			return nil, err
		}

		metric = &models.Metrics{
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

		metric = &models.Metrics{
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
