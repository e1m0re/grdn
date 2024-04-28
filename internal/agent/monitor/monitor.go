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
	"github.com/e1m0re/grdn/internal/storage"
)

type MetricsState struct {
	Gauges   map[storage.GaugeName]storage.GaugeDateType
	Counters map[storage.CounterName]storage.CounterDateType
}

type MetricsMonitor struct {
	mx   sync.RWMutex
	data MetricsState
}

func NewMetricsMonitor() *MetricsMonitor {
	return &MetricsMonitor{
		data: MetricsState{
			Gauges:   make(map[storage.GaugeName]storage.GaugeDateType),
			Counters: make(map[storage.CounterName]storage.CounterDateType),
		},
	}
}

func (m *MetricsMonitor) UpdateData() {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.data.Counters[storage.PollCount]++

	r := rand.New(rand.NewSource(m.data.Counters[storage.PollCount]))
	m.data.Gauges[storage.RandomValue] = r.Float64()

	// memory metrics
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)
	m.data.Gauges[storage.Alloc] = storage.GaugeDateType(rtm.Alloc)
	m.data.Gauges[storage.BuckHashSys] = storage.GaugeDateType(rtm.BuckHashSys)
	m.data.Gauges[storage.Frees] = storage.GaugeDateType(rtm.Frees)
	m.data.Gauges[storage.GCCPUFraction] = rtm.GCCPUFraction
	m.data.Gauges[storage.GCSys] = storage.GaugeDateType(rtm.GCSys)
	m.data.Gauges[storage.HeapAlloc] = storage.GaugeDateType(rtm.HeapAlloc)
	m.data.Gauges[storage.HeapIdle] = storage.GaugeDateType(rtm.HeapIdle)
	m.data.Gauges[storage.HeapInuse] = storage.GaugeDateType(rtm.HeapInuse)
	m.data.Gauges[storage.HeapObjects] = storage.GaugeDateType(rtm.HeapObjects)
	m.data.Gauges[storage.HeapReleased] = storage.GaugeDateType(rtm.HeapReleased)
	m.data.Gauges[storage.HeapSys] = storage.GaugeDateType(rtm.HeapSys)
	m.data.Gauges[storage.LastGC] = storage.GaugeDateType(rtm.LastGC)
	m.data.Gauges[storage.Lookups] = storage.GaugeDateType(rtm.Lookups)
	m.data.Gauges[storage.MCacheInuse] = storage.GaugeDateType(rtm.MCacheInuse)
	m.data.Gauges[storage.MCacheSys] = storage.GaugeDateType(rtm.MCacheSys)
	m.data.Gauges[storage.MSpanInuse] = storage.GaugeDateType(rtm.MSpanInuse)
	m.data.Gauges[storage.MSpanSys] = storage.GaugeDateType(rtm.MSpanSys)
	m.data.Gauges[storage.Mallocs] = storage.GaugeDateType(rtm.Mallocs)
	m.data.Gauges[storage.NextGC] = storage.GaugeDateType(rtm.NextGC)
	m.data.Gauges[storage.NumForcedGC] = storage.GaugeDateType(rtm.NumForcedGC)
	m.data.Gauges[storage.NumGC] = storage.GaugeDateType(rtm.NumGC)
	m.data.Gauges[storage.OtherSys] = storage.GaugeDateType(rtm.OtherSys)
	m.data.Gauges[storage.StackInuse] = storage.GaugeDateType(rtm.StackInuse)
	m.data.Gauges[storage.StackSys] = storage.GaugeDateType(rtm.StackSys)
	m.data.Gauges[storage.PauseTotalNs] = storage.GaugeDateType(rtm.PauseTotalNs)
	m.data.Gauges[storage.Sys] = storage.GaugeDateType(rtm.Sys)
	m.data.Gauges[storage.TotalAlloc] = storage.GaugeDateType(rtm.TotalAlloc)
}

func (m *MetricsMonitor) UpdateGOPS(ctx context.Context) {
	m.mx.Lock()
	defer m.mx.Unlock()

	memoryInfo, _ := mem.VirtualMemoryWithContext(ctx)

	m.data.Gauges["TotalMemory"] = storage.GaugeDateType(memoryInfo.Total)
	m.data.Gauges["FreeMemory"] = storage.GaugeDateType(memoryInfo.Free)

	percents, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		panic(err)
	}

	for idx, percent := range percents {
		m.data.Gauges[fmt.Sprintf("CPUutilization%d", idx)] = percent
	}
}

func (m *MetricsMonitor) GetMetricsList() models.MetricsList {
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
