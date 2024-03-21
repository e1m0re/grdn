package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
)

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

func TestMemStorage_UpdateMetricValue(t *testing.T) {
	type fields struct {
		Gauges   map[GaugeName]GaugeDateType
		Counters map[CounterName]CounterDateType
	}
	type args struct {
		mType  models.MetricsType
		mName  string
		mValue string
	}
	type want struct {
		store fields
		err   string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "test success update gauge metric",
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
			args: args{
				mType:  "gauge",
				mName:  "metric",
				mValue: "123.123",
			},
			want: want{
				store: fields{
					Gauges: map[GaugeName]GaugeDateType{
						"metric": 123.123,
					},
					Counters: make(map[CounterName]CounterDateType),
				},
				err: "",
			},
		},
		{
			name: "test success update counter metric",
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
			args: args{
				mType:  "counter",
				mName:  "metric",
				mValue: "123",
			},
			want: want{
				store: fields{
					Gauges: make(map[GaugeName]GaugeDateType),
					Counters: map[CounterName]CounterDateType{
						"metric": 123,
					},
				},
				err: "",
			},
		},
		{
			name: "test update invalid type metric",
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
			args: args{
				mType:  "new_counter",
				mName:  "metric",
				mValue: "123",
			},
			want: want{
				store: fields{
					Gauges:   make(map[GaugeName]GaugeDateType),
					Counters: make(map[CounterName]CounterDateType),
				},
				err: "unknown metric type",
			},
		},
		{
			name: "test update gauge metric with invalid value",
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
			args: args{
				mType:  "gauge",
				mName:  "metric",
				mValue: "qwerty",
			},
			want: want{
				store: fields{
					Gauges:   make(map[GaugeName]GaugeDateType),
					Counters: make(map[CounterName]CounterDateType),
				},
				err: "strconv.ParseFloat: parsing \"qwerty\": invalid syntax",
			},
		},
		{
			name: "test update counter metric with invalid value",
			fields: fields{
				Gauges:   make(map[GaugeName]GaugeDateType),
				Counters: make(map[CounterName]CounterDateType),
			},
			args: args{
				mType:  "counter",
				mName:  "metric",
				mValue: "qwerty",
			},
			want: want{
				store: fields{
					Gauges:   make(map[GaugeName]GaugeDateType),
					Counters: make(map[CounterName]CounterDateType),
				},
				err: "strconv.ParseInt: parsing \"qwerty\": invalid syntax",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}

			err := s.UpdateMetricValue(tt.args.mType, tt.args.mName, tt.args.mValue)
			if tt.want.err == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.want.store.Gauges, s.Gauges)
				assert.Equal(t, tt.want.store.Counters, s.Counters)

				return
			}

			require.Error(t, err)
			assert.EqualError(t, err, tt.want.err)
			assert.Equal(t, tt.want.store.Gauges, s.Gauges)
			assert.Equal(t, tt.want.store.Counters, s.Counters)
		})
	}
}
