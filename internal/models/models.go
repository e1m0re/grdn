// Package models contains common application models.
package models

import (
	"fmt"
	"strconv"
)

type MetricsType = string

const (
	GaugeType   = MetricsType("gauge")
	CounterType = MetricsType("counter")
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

type Metric struct {
	Value *float64    `json:"value,omitempty" db:"value"`
	Delta *int64      `json:"delta,omitempty" db:"delta"`
	MType MetricsType `json:"type" db:"type"`
	ID    string      `json:"id" db:"name"`
}

func (m *Metric) ValueToString() string {
	switch m.MType {
	case GaugeType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case CounterType:
		return fmt.Sprintf("%d", *m.Delta)
	default:
		return ""
	}
}

func (m *Metric) String() string {
	return fmt.Sprintf("%s: %s", m.ID, m.ValueToString())
}

func (m *Metric) ValueFromString(str string) error {
	switch m.MType {
	case GaugeType:
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}

		m.Value = &value
	case CounterType:
		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		m.Delta = &value
	}

	return nil
}

type MetricsList []*Metric
