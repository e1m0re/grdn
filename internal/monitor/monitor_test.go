package monitor

import (
	"github.com/e1m0re/grdn/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsMonitor_GetData(t *testing.T) {
	type fields struct {
		data MetricsState
	}

	tests := []struct {
		name   string
		fields fields
		want   []GlobalMetricsList
	}{
		{
			name: "test get data with ALL types metrics",
			fields: fields{data: MetricsState{
				Gauges:   map[storage.GaugeName]storage.GaugeDateType{"Alloc": 123.123},
				Counters: map[storage.CounterName]storage.CounterDateType{"PollCount": 10},
			}},
			want: []GlobalMetricsList{
				{
					MType:  storage.GaugeType,
					MName:  "Alloc",
					MValue: "123.123",
				},
				{
					MType:  storage.CounterType,
					MName:  "PollCount",
					MValue: "10",
				},
			},
		},
		{
			name: "test get data without gauge types metrics",
			fields: fields{data: MetricsState{
				Gauges:   make(map[storage.GaugeName]storage.GaugeDateType),
				Counters: map[storage.CounterName]storage.CounterDateType{"PollCount": 10},
			}},
			want: []GlobalMetricsList{
				{
					MType:  storage.CounterType,
					MName:  "PollCount",
					MValue: "10",
				},
			},
		},
		{
			name: "test get data without counter types metrics",
			fields: fields{data: MetricsState{
				Gauges:   map[storage.GaugeName]storage.GaugeDateType{"Alloc": 123.123},
				Counters: make(map[storage.CounterName]storage.CounterDateType),
			}},
			want: []GlobalMetricsList{
				{
					MType:  storage.GaugeType,
					MName:  "Alloc",
					MValue: "123.123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricsMonitor{
				data: tt.fields.data,
			}
			for idx, row := range m.GetData() {
				assert.Equal(t, tt.want[idx], row)
			}
		})
	}
}

func TestMetricsMonitor_UpdateData(t *testing.T) {
	type fields struct {
		data MetricsState
	}

	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{
			name: "test 1",
			fields: fields{data: MetricsState{
				Gauges:   make(map[storage.GaugeName]storage.GaugeDateType),
				Counters: make(map[storage.CounterName]storage.CounterDateType),
			}},
			want: fields{data: MetricsState{
				Gauges: map[storage.GaugeName]storage.GaugeDateType{
					storage.RandomValue:   0,
					storage.Alloc:         0,
					storage.BuckHashSys:   0,
					storage.Frees:         0,
					storage.GCCPUFraction: 0,
					storage.GCSys:         0,
					storage.HeapAlloc:     0,
					storage.HeapIdle:      0,
					storage.HeapInuse:     0,
					storage.HeapObjects:   0,
					storage.HeapReleased:  0,
					storage.HeapSys:       0,
					storage.LastGC:        0,
					storage.Lookups:       0,
					storage.MCacheInuse:   0,
					storage.MCacheSys:     0,
					storage.MSpanInuse:    0,
					storage.MSpanSys:      0,
					storage.Mallocs:       0,
					storage.NextGC:        0,
					storage.NumForcedGC:   0,
					storage.NumGC:         0,
					storage.OtherSys:      0,
					storage.StackInuse:    0,
					storage.StackSys:      0,
					storage.PauseTotalNs:  0,
					storage.Sys:           0,
					storage.TotalAlloc:    0,
				},
				Counters: map[storage.CounterName]storage.CounterDateType{
					storage.PollCount: 1,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricsMonitor{
				data: tt.fields.data,
			}
			m.UpdateData()

			for key := range tt.want.data.Gauges {
				_, ok := m.data.Gauges[key]
				assert.True(t, ok)
			}

			for key := range tt.want.data.Counters {
				_, ok := m.data.Counters[key]
				assert.True(t, ok)
			}
		})
	}
}

func TestNewMetricsMonitor(t *testing.T) {
	tests := []struct {
		name string
		want MetricsMonitor
	}{
		{
			name: "test MetricsMonitor constructor",
			want: MetricsMonitor{
				data: MetricsState{
					Gauges:   make(map[storage.GaugeName]storage.GaugeDateType),
					Counters: make(map[storage.CounterName]storage.CounterDateType),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor1 := NewMetricsMonitor()
			assert.Equal(t, &tt.want, monitor1)
		})
	}
}
