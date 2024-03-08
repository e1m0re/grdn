package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCounterName(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test PollCount",
			args: args{value: "PollCount"},
			want: true,
		},
		{
			name: "test empty string",
			args: args{value: ""},
			want: false,
		},
		{
			name: "test Invalid name",
			args: args{value: "Invalid name"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidCounterName(tt.args.value))
		})
	}
}

func TestIsValidGuageName(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test Alloc",
			args: args{value: "Alloc"},
			want: true,
		},
		{
			name: "test BuckHashSys",
			args: args{value: "BuckHashSys"},
			want: true,
		},
		{
			name: "test empty string",
			args: args{value: ""},
			want: false,
		},
		{
			name: "test Invalid name",
			args: args{value: "Invalid name"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidGaugeName(tt.args.value))
		})
	}
}

func TestIsValidMetricsType(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test Guage => true",
			args: args{
				value: GaugeType,
			},
			want: true,
		},
		{
			name: "test Counter => true",
			args: args{
				value: CounterType,
			},
			want: true,
		},
		{
			name: "test empty string",
			args: args{
				value: "",
			},
			want: false,
		},
		{
			name: "test invalid type",
			args: args{
				value: "Invalid type",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidMetricsType(tt.args.value))
		})
	}
}

func TestMemStorage_UpdateCounterMetric(t *testing.T) {
	type fields struct {
		Gauges   map[GaugeName]GaugeDateType
		Counters map[CounterName]CounterDateType
	}
	type args struct {
		name  CounterName
		value CounterDateType
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   fields
	}{
		{
			name: "test update PollCount",
			args: args{
				name:  "PollCount",
				value: 200,
			},
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: map[CounterName]CounterDateType{"PollCount": 100},
			},
			want: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: map[CounterName]CounterDateType{"PollCount": 300},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.UpdateCounterMetric(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.Counters[tt.args.name], s.Counters[tt.args.name])
		})
	}
}

func TestMemStorage_UpdateGuageMetric(t *testing.T) {
	type fields struct {
		Gauges   map[GaugeName]GaugeDateType
		Counters map[CounterName]CounterDateType
	}
	type args struct {
		name  GaugeName
		value GaugeDateType
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   fields
	}{
		{
			name: "test update Alloc",
			args: args{
				name:  "Alloc",
				value: 200,
			},
			fields: fields{
				Gauges:   map[GaugeName]GaugeDateType{"Alloc": 100},
				Counters: make(map[CounterName]CounterDateType),
			},
			want: fields{
				Gauges:   map[GaugeName]GaugeDateType{"Alloc": 200},
				Counters: make(map[CounterName]CounterDateType),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.UpdateGaugeMetric(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.Counters[tt.args.name], s.Counters[tt.args.name])
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "test constructor",
			want: &MemStorage{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage()
			assert.Equal(t, tt.want, storage)
			assert.Equal(t, *tt.want, *storage)
		})
	}
}

func TestNewMetricsState(t *testing.T) {
	tests := []struct {
		name string
		want *MetricsState
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricsState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricsState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_IsValidMetricName(t *testing.T) {
	type fields struct {
		Gauges   map[GaugeName]GaugeDateType
		Counters map[CounterName]CounterDateType
	}
	type args struct {
		mType MetricsType
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test check gauge type",
			args: args{
				mType: GaugeType,
				value: Alloc,
			},
			want: true,
		},
		{
			name: "test check counter type",
			args: args{
				mType: CounterType,
				value: "PollCount",
			},
			want: true,
		},
		{
			name: "test empty type",
			args: args{
				mType: "",
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown type",
			args: args{
				mType: "Unknown",
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown value for guage type",
			args: args{
				mType: GaugeType,
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown value for Counter type",
			args: args{
				mType: CounterType,
				value: "PollCount1",
			},
			want: false,
		},
		{
			name: "test empty type",
			args: args{
				mType: CounterType,
				value: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			assert.Equalf(t, tt.want, s.IsValidMetricName(tt.args.mType, tt.args.value), "IsValidMetricName(%v, %v)", tt.args.mType, tt.args.value)
		})
	}
}
