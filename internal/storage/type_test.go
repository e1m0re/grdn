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
			assert.Equal(t, tt.want, IsValidGuageName(tt.args.value))
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
				value: GuageType,
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
		Guages   map[GuageName]GuageDateType
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
				Guages:   make(map[GuageName]GuageDateType),
				Counters: map[CounterName]CounterDateType{"PollCount": 100},
			},
			want: fields{
				Guages:   make(map[GuageName]GuageDateType),
				Counters: map[CounterName]CounterDateType{"PollCount": 300},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Guages:   tt.fields.Guages,
				Counters: tt.fields.Counters,
			}
			s.UpdateCounterMetric(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.Counters[tt.args.name], s.Counters[tt.args.name])
		})
	}
}

func TestMemStorage_UpdateGuageMetric(t *testing.T) {
	type fields struct {
		Guages   map[GuageName]GuageDateType
		Counters map[CounterName]CounterDateType
	}
	type args struct {
		name  GuageName
		value GuageDateType
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
				Guages:   map[GuageName]GuageDateType{"Alloc": 100},
				Counters: make(map[CounterName]CounterDateType),
			},
			want: fields{
				Guages:   map[GuageName]GuageDateType{"Alloc": 200},
				Counters: make(map[CounterName]CounterDateType),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Guages:   tt.fields.Guages,
				Counters: tt.fields.Counters,
			}
			s.UpdateGuageMetric(tt.args.name, tt.args.value)
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
				Guages:   make(map[GuageName]GuageDateType),
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
