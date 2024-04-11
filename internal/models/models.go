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

type Metric struct {
	ID    string      `json:"id" db:"name"`
	MType MetricsType `json:"type" db:"type"`
	Delta *int64      `json:"delta,omitempty" db:"delta"`
	Value *float64    `json:"value,omitempty" db:"value"`
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
