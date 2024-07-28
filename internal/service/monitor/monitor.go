// Package monitor implements business logic of clients application.
package monitor

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/e1m0re/grdn/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Monitor
type Monitor interface {
	// UpdateData collects main metrics data and updates local state.
	UpdateData()
	// UpdateGOPS collects add metrics data and updates local state.
	UpdateGOPS(ctx context.Context) error
	// GetMetricsList returns all metrics.
	GetMetricsList() models.MetricsList
}

type MetricsState struct {
	Gauges   map[models.GaugeName]models.GaugeDateType
	Counters map[models.CounterName]models.CounterDateType
}

type monitor struct {
	data MetricsState
	mx   sync.RWMutex
}

// NewMonitor is MetricsMonitor constructor.
func NewMonitor() Monitor {
	return &monitor{
		data: MetricsState{
			Gauges:   make(map[models.GaugeName]models.GaugeDateType),
			Counters: make(map[models.CounterName]models.CounterDateType),
		},
	}
}

// UpdateData collects main metrics data and updates local state.
func (m *monitor) UpdateData() {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.data.Counters[models.PollCount]++

	r := rand.New(rand.NewSource(m.data.Counters[models.PollCount]))
	m.data.Gauges[models.RandomValue] = r.Float64()

	// memory metrics
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)
	m.data.Gauges[models.Alloc] = models.GaugeDateType(rtm.Alloc)
	m.data.Gauges[models.BuckHashSys] = models.GaugeDateType(rtm.BuckHashSys)
	m.data.Gauges[models.Frees] = models.GaugeDateType(rtm.Frees)
	m.data.Gauges[models.GCCPUFraction] = rtm.GCCPUFraction
	m.data.Gauges[models.GCSys] = models.GaugeDateType(rtm.GCSys)
	m.data.Gauges[models.HeapAlloc] = models.GaugeDateType(rtm.HeapAlloc)
	m.data.Gauges[models.HeapIdle] = models.GaugeDateType(rtm.HeapIdle)
	m.data.Gauges[models.HeapInuse] = models.GaugeDateType(rtm.HeapInuse)
	m.data.Gauges[models.HeapObjects] = models.GaugeDateType(rtm.HeapObjects)
	m.data.Gauges[models.HeapReleased] = models.GaugeDateType(rtm.HeapReleased)
	m.data.Gauges[models.HeapSys] = models.GaugeDateType(rtm.HeapSys)
	m.data.Gauges[models.LastGC] = models.GaugeDateType(rtm.LastGC)
	m.data.Gauges[models.Lookups] = models.GaugeDateType(rtm.Lookups)
	m.data.Gauges[models.MCacheInuse] = models.GaugeDateType(rtm.MCacheInuse)
	m.data.Gauges[models.MCacheSys] = models.GaugeDateType(rtm.MCacheSys)
	m.data.Gauges[models.MSpanInuse] = models.GaugeDateType(rtm.MSpanInuse)
	m.data.Gauges[models.MSpanSys] = models.GaugeDateType(rtm.MSpanSys)
	m.data.Gauges[models.Mallocs] = models.GaugeDateType(rtm.Mallocs)
	m.data.Gauges[models.NextGC] = models.GaugeDateType(rtm.NextGC)
	m.data.Gauges[models.NumForcedGC] = models.GaugeDateType(rtm.NumForcedGC)
	m.data.Gauges[models.NumGC] = models.GaugeDateType(rtm.NumGC)
	m.data.Gauges[models.OtherSys] = models.GaugeDateType(rtm.OtherSys)
	m.data.Gauges[models.StackInuse] = models.GaugeDateType(rtm.StackInuse)
	m.data.Gauges[models.StackSys] = models.GaugeDateType(rtm.StackSys)
	m.data.Gauges[models.PauseTotalNs] = models.GaugeDateType(rtm.PauseTotalNs)
	m.data.Gauges[models.Sys] = models.GaugeDateType(rtm.Sys)
	m.data.Gauges[models.TotalAlloc] = models.GaugeDateType(rtm.TotalAlloc)
}

// UpdateGOPS collects add metrics data and updates local state.
func (m *monitor) UpdateGOPS(ctx context.Context) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	memoryInfo, _ := mem.VirtualMemoryWithContext(ctx)

	m.data.Gauges[models.TotalMemory] = models.GaugeDateType(memoryInfo.Total)
	m.data.Gauges[models.FreeMemory] = models.GaugeDateType(memoryInfo.Free)

	percents, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return err
	}

	for idx, percent := range percents {
		m.data.Gauges[fmt.Sprintf("CPUutilization%d", idx)] = percent
	}

	return nil
}

// GetMetricsList returns all metrics.
func (m *monitor) GetMetricsList() models.MetricsList {
	m.mx.RLock()
	defer m.mx.RUnlock()

	result := make(models.MetricsList, 0)
	for key, value := range m.data.Gauges {
		x := value
		result = append(result, &models.Metric{
			ID:    key,
			MType: models.GaugeType,
			Value: &x,
		})
	}

	for key, value := range m.data.Counters {
		x := value
		result = append(result, &models.Metric{
			ID:    key,
			MType: models.CounterType,
			Delta: &x,
		})
	}

	return result
}
