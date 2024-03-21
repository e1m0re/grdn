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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) ValueToString() string {
	switch m.MType {
	case GaugeType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case CounterType:
		return fmt.Sprintf("%d", *m.Delta)
	default:
		return ""
	}
}

type MetricsList []*Metrics
