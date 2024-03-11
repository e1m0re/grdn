package storage

import (
	"errors"
	"fmt"
	"strconv"
)

type MetricsType = string

const (
	GaugeType   = MetricsType("gauge")
	CounterType = MetricsType("counter")
)

func IsValidMetricsType(value string) bool {
	switch value {
	case GaugeType, CounterType:
		return true
	default:
		return false
	}
}

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

func IsValidGaugeName(value string) bool {
	switch value {
	case Alloc, BuckHashSys, Frees, GCCPUFraction, GCSys, HeapAlloc, HeapIdle, HeapInuse, HeapObjects, HeapReleased,
		HeapSys, LastGC, Lookups, MCacheInuse, MCacheSys, MSpanInuse, MSpanSys, Mallocs, NextGC, NumForcedGC, NumGC,
		OtherSys, PauseTotalNs, StackInuse, StackSys, Sys, TotalAlloc, RandomValue:
		return true
	default:
		return false
	}
}

type CounterDateType = int64
type CounterName = string

const (
	PollCount = CounterName("PollCount")
)

func IsValidCounterName(value string) bool {
	switch value {
	case PollCount:
		return true
	}

	return false
}

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

func (s *MemStorage) GetMetricValue(mType MetricsType, mName string) (string, error) {
	switch mType {
	case GaugeType:
		if value, ok := s.Gauges[mName]; ok {
			return strconv.FormatFloat(value, 'f', -1, 64), nil
		}
		return "", errors.New("unknown metric")
	case CounterType:
		if value, ok := s.Counters[mName]; ok {
			return fmt.Sprintf("%d", value), nil
		}
		return "", errors.New("unknown metric")
	default:
		return "", errors.New("unknown metric")
	}
}

type MetricsState struct {
	Gauges   map[GaugeName]GaugeDateType
	Counters map[CounterName]CounterDateType
}

func NewMetricsState() *MetricsState {
	return &MetricsState{
		Gauges:   make(map[GaugeName]GaugeDateType),
		Counters: make(map[CounterName]CounterDateType),
	}
}
