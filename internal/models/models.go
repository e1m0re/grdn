package models

import (
	"fmt"
	"strconv"
)

type MetricType = string
type MetricName = string

const (
	GaugeType   = MetricType("gauge")
	CounterType = MetricType("counter")
)

type Metric struct {
	Value *float64   `json:"value,omitempty" db:"value"`
	Delta *int64     `json:"delta,omitempty" db:"delta"`
	MType MetricType `json:"type" db:"type"`
	ID    MetricName `json:"id" db:"name"`
}

type MetricsList []*Metric

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
