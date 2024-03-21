package storage

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/e1m0re/grdn/internal/models"
)

type GaugeDateType = float64
type GaugeName = string

const (
	Alloc         = GaugeName("Alloc")
	BuckHashSys   = GaugeName("BuckHashSys")
	Frees         = GaugeName("Frees")
	GCCPUFraction = GaugeName("GCCPUFraction")
	GCSys         = GaugeName("GCSys")
	HeapAlloc     = GaugeName("HeapAlloc")
	HeapIdle      = GaugeName("HeapIdle")
	HeapInuse     = GaugeName("HeapInuse")
	HeapObjects   = GaugeName("HeapObjects")
	HeapReleased  = GaugeName("HeapReleased")
	HeapSys       = GaugeName("HeapSys")
	LastGC        = GaugeName("LastGC")
	Lookups       = GaugeName("Lookups")
	MCacheInuse   = GaugeName("MCacheInuse")
	MCacheSys     = GaugeName("MCacheSys")
	MSpanInuse    = GaugeName("MSpanInuse")
	MSpanSys      = GaugeName("MSpanSys")
	Mallocs       = GaugeName("Mallocs")
	NextGC        = GaugeName("NextGC")
	NumForcedGC   = GaugeName("NumForcedGC")
	NumGC         = GaugeName("NumGC")
	OtherSys      = GaugeName("OtherSys")
	PauseTotalNs  = GaugeName("PauseTotalNs")
	StackInuse    = GaugeName("StackInuse")
	StackSys      = GaugeName("StackSys")
	Sys           = GaugeName("Sys")
	TotalAlloc    = GaugeName("TotalAlloc")
	RandomValue   = GaugeName("RandomValue")
)

type CounterDateType = int64
type CounterName = string

const (
	PollCount = CounterName("PollCount")
)

type MemStorage struct {
	Gauges   map[GaugeName]GaugeDateType
	Counters map[CounterName]CounterDateType
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[GaugeName]GaugeDateType),
		Counters: make(map[CounterName]CounterDateType),
	}
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

	return nil
}

func (s *MemStorage) GetAllMetrics() []string {
	var result []string
	for key, value := range s.Gauges {
		result = append(result, fmt.Sprintf("%s: %s", key, strconv.FormatFloat(value, 'f', -1, 64)))
	}

	for key, value := range s.Counters {
		result = append(result, fmt.Sprintf("%s: %v", key, value))
	}

	return result
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
