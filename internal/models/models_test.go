package models

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMetric_ValueFromString(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		MType MetricType
		ID    MetricName
	}
	type args struct {
		str string
	}
	type want struct {
		err    error
		metric Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Failed parse int value",
			fields: fields{
				MType: CounterType,
				ID:    "metric 1",
			},
			args: args{
				str: "10a0",
			},
			want: want{
				err: &strconv.NumError{
					Func: "ParseInt",
					Num:  "10a0",
					Err:  strconv.ErrSyntax,
				},
				metric: Metric{
					MType: CounterType,
					ID:    "metric 1",
				},
			},
		},
		{
			name: "Failed parse float value",
			fields: fields{
				MType: GaugeType,
				ID:    "metric 1",
			},
			args: args{
				str: "10a0.01",
			},
			want: want{
				err: &strconv.NumError{
					Func: "ParseFloat",
					Num:  "10a0.01",
					Err:  strconv.ErrSyntax,
				},
				metric: Metric{
					MType: GaugeType,
					ID:    "metric 1",
				},
			},
		},
		{
			name: "Successfully case (counter)",
			fields: fields{
				MType: CounterType,
				ID:    "metric 1",
			},
			args: args{
				str: "100",
			},
			want: want{
				err: nil,
				metric: Metric{
					MType: CounterType,
					ID:    "metric 1",
					Delta: &d,
				},
			},
		},
		{
			name: "Successfully case (gauge)",
			fields: fields{
				MType: GaugeType,
				ID:    "metric 1",
			},
			args: args{
				str: "100.1",
			},
			want: want{
				err: nil,
				metric: Metric{
					MType: GaugeType,
					ID:    "metric 1",
					Value: &v,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := Metric{
				MType: test.fields.MType,
				ID:    test.fields.ID,
			}
			err := m.ValueFromString(test.args.str)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.metric, m)
		})
	}
}

func TestMetric_ValueToString(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		Value *float64
		Delta *int64
		MType MetricType
		ID    MetricName
	}
	type want struct {
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name:   "unknown type",
			fields: fields{},
			want: want{
				result: "",
			},
		},
		{
			name: "successfully case (counter)",
			fields: fields{
				Delta: &d,
				MType: CounterType,
			},
			want: want{
				result: "100",
			},
		},
		{
			name: "successfully case (gauge)",
			fields: fields{
				Value: &v,
				MType: GaugeType,
			},
			want: want{
				result: "100.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				Value: tt.fields.Value,
				Delta: tt.fields.Delta,
				MType: tt.fields.MType,
				ID:    tt.fields.ID,
			}
			assert.Equalf(t, tt.want.result, m.ValueToString(), "ValueToString()")
		})
	}
}

func TestMetric_String(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		Value *float64
		Delta *int64
		MType MetricType
		ID    MetricName
	}
	type want struct {
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "stringer counter",
			fields: fields{
				Delta: &d,
				MType: CounterType,
				ID:    "metric 1",
			},
			want: want{result: "metric 1: 100"},
		},
		{
			name: "gauge counter",
			fields: fields{
				Value: &v,
				MType: GaugeType,
				ID:    "metric 1",
			},
			want: want{result: "metric 1: 100.1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &Metric{
				Value: test.fields.Value,
				Delta: test.fields.Delta,
				MType: test.fields.MType,
				ID:    test.fields.ID,
			}
			assert.Equalf(t, test.want.result, m.String(), "String()")
		})
	}
}
