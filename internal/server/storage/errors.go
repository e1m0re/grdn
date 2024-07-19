package storage

import "errors"

var (
	ErrUnknownMetric     = errors.New("unknown metric")
	ErrUnknownMetricType = errors.New("unknown metric type")
)
