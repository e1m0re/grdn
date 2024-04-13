package storage

import (
	"context"

	"github.com/e1m0re/grdn/internal/models"
)

type Interface interface {
	Close() error
	DumpStorageToFile() error
	GetAllMetrics(ctx context.Context) ([]string, error)
	GetMetric(ctx context.Context, mType models.MetricsType, mName string) (metric *models.Metric, err error)
	LoadStorageFromFile() error
	Ping(ctx context.Context) error
	UpdateMetric(ctx context.Context, metric models.Metric) error
	UpdateMetrics(ctx context.Context, metrics models.MetricsList) error
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

type CounterDateType = int64
type CounterName = string

const (
	PollCount = CounterName("PollCount")
)
