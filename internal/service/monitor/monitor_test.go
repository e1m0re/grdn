package monitor

import (
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
			m := &MetricsMonitor{
				data: tt.fields.data,
			}

			result := m.GetMetricsList()

			assert.Equal(t, tt.want, result)
		})
	}
}
