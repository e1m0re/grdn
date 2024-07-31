package monitor

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/e1m0re/grdn/internal/models"
)

func TestMetricsMonitor_GetMetricsList(t *testing.T) {
	v1 := 8.07
	d1 := int64(1984)
	type fields struct {
		data MetricsState
	}
	tests := []struct {
		name   string
		fields fields
		want   models.MetricsList
	}{
		{
			name: "test empty store",
			fields: fields{
				data: MetricsState{
					make(map[models.GaugeName]models.GaugeDateType),
					make(map[models.CounterName]models.CounterDateType),
				},
			},
			want: make(models.MetricsList, 0),
		},
		{
			name: "test not empty store",
			fields: fields{
				data: MetricsState{
					Gauges: map[models.GaugeName]models.GaugeDateType{
						"metric1": v1,
					},
					Counters: map[models.CounterName]models.CounterDateType{
						"metric2": d1,
					},
				},
			},
			want: models.MetricsList{
				&models.Metric{
					ID:    "metric1",
					MType: "gauge",
					Delta: nil,
					Value: &v1,
				},
				&models.Metric{
					ID:    "metric2",
					MType: "counter",
					Delta: &d1,
					Value: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &monitor{
				data: tt.fields.data,
			}

			result := m.GetMetricsList()

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMetricsMonitor_UpdateGOPS(t *testing.T) {
	type fields struct {
		data MetricsState
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err          error
		metricsNames []models.MetricName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Successfully case",
			fields: fields{
				data: MetricsState{
					Gauges:   make(map[models.GaugeName]models.GaugeDateType),
					Counters: make(map[models.CounterName]models.CounterDateType),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
				metricsNames: []models.MetricName{
					models.TotalMemory,
					models.FreeMemory,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &monitor{
				data: test.fields.data,
			}
			err := m.UpdateGOPS(test.args.ctx)
			assert.Equal(t, test.want.err, err)
			fmt.Printf("------------------------ %v", m)
			for _, name := range test.want.metricsNames {
				assert.Contains(t, m.data.Gauges, name)
			}
			for i := 0; i < runtime.NumCPU(); i++ {
				assert.Contains(t, m.data.Gauges, fmt.Sprintf("CPUutilization%d", i))
			}
		})
	}
}

func TestMetricsMonitor_UpdateData(t *testing.T) {
	type fields struct {
		data MetricsState
	}
	type want struct {
		metricsGaugeNames   []models.MetricName
		metricsCounterNames []models.MetricName
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Successfully case",
			fields: fields{
				data: MetricsState{
					Gauges:   make(map[models.GaugeName]models.GaugeDateType),
					Counters: make(map[models.CounterName]models.CounterDateType),
				},
			},
			want: want{
				metricsGaugeNames:   models.MetricsGaugeNamesList,
				metricsCounterNames: models.MetricsCounterNamesList,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &monitor{
				data: test.fields.data,
			}
			m.UpdateData()
			for _, name := range test.want.metricsGaugeNames {
				assert.Contains(t, m.data.Gauges, name)
			}
			for _, name := range test.want.metricsCounterNames {
				assert.Contains(t, m.data.Counters, name)
			}
		})
	}
}

func TestNewMetricsMonitor(t *testing.T) {
	tests := []struct {
		want Monitor
		name string
	}{
		{
			name: "Successfully case",
			want: &monitor{
				data: MetricsState{
					Gauges:   make(map[models.GaugeName]models.GaugeDateType),
					Counters: make(map[models.CounterName]models.CounterDateType),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := NewMonitor()
			assert.Equalf(t, test.want, m, "NewMetricsMonitor()")
			assert.Implements(t, (*Monitor)(nil), m)
		})
	}
}
