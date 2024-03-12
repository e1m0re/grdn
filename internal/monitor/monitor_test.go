package monitor

import (
	"reflect"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricsMonitor{
				data: tt.fields.data,
			}
			if got := m.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricsMonitor{
				data: tt.fields.data,
			}
			m.UpdateData()
		})
	}
}

func TestNewMetricsMonitor(t *testing.T) {
	tests := []struct {
		name string
		want *MetricsMonitor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricsMonitor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricsMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}
