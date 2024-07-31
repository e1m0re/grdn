package models

type CounterDateType = int64
type CounterName = MetricName

const (
	PollCount = CounterName("PollCount")
)

var MetricsCounterNamesList = []MetricName{
	PollCount,
}
