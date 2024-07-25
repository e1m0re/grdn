package memory

import (
	"context"
	"fmt"
	"github.com/e1m0re/grdn/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestStore_Clear(t *testing.T) {
	type fields struct {
		metrics map[string]models.Metric
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err          error
		metricsCount int
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
				metrics: map[string]models.Metric{"metric1": {ID: "metric1"}},
			},
			args: args{ctx: nil},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{
				metrics: test.fields.metrics,
			}

			err := s.Clear(test.args.ctx)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.metricsCount, len(s.metrics))
		})
	}
}

func TestStore_genMetricKey(t *testing.T) {
	type args struct {
		m models.MetricName
		t models.MetricType
	}
	type want struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Empty metric name",
			args: args{
				m: "",
				t: models.CounterType,
			},
			want: want{
				key: "886bb73b3156b0aa24aac99d2de0b238",
			},
		},
		{
			name: "Empty metric type",
			args: args{
				m: "metric 1",
				t: "",
			},
			want: want{
				key: "65dba30d4a32b317a4ffebd2f077a49b",
			},
		},
		{
			name: "Successfully case for counter metric",
			args: args{
				m: "metric 1",
				t: models.CounterType,
			},
			want: want{
				key: "c917499366554154b438bb8259183adc",
			},
		},
		{
			name: "Successfully case for gauge metric",
			args: args{
				m: "metric 1",
				t: models.GaugeType,
			},
			want: want{
				key: "9e646d6d5855fbdadcaab202f1748505",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{}
			got := s.genMetricKey(test.args.m, test.args.t)
			assert.Equal(t, test.want.key, got)
		})
	}
}

func TestStore_Close(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		want want
		name string
	}{
		{
			name: "Successfully case",
			want: want{err: nil},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{}
			err := s.Close()
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		args args
		want want
		name string
	}{
		{
			name: "Successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{err: nil},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{}
			err := s.Ping(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_GetMetric(t *testing.T) {
	d := int64(100)
	type fields struct {
		metrics map[string]models.Metric
	}
	type args struct {
		ctx   context.Context
		mType models.MetricType
		mName string
	}
	type want struct {
		metric *models.Metric
		err    error
	}
	tests := []struct {
		fields fields
		want   want
		args   args
		name   string
	}{
		{
			name: "Unknown metric",
			fields: fields{
				metrics: make(map[string]models.Metric),
			},
			args: args{
				ctx:   context.Background(),
				mType: models.CounterType,
				mName: "metric 1",
			},
			want: want{
				metric: nil,
				err:    nil,
			},
		},
		{
			name: "Successfully case",
			fields: fields{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
				},
			},
			args: args{
				ctx:   context.Background(),
				mType: models.CounterType,
				mName: "metric 1",
			},
			want: want{
				metric: &models.Metric{
					Value: nil,
					Delta: &d,
					MType: models.CounterType,
					ID:    "metric 1",
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{
				metrics: test.fields.metrics,
			}
			got, err := s.GetMetric(test.args.ctx, test.args.mType, test.args.mName)
			assert.Equal(t, test.want.metric, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_GetAllMetrics(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		metrics map[string]models.Metric
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err     error
		metrics models.MetricsList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Empty list",
			fields: fields{
				metrics: make(map[string]models.Metric),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				metrics: make(models.MetricsList, 0),
				err:     nil,
			},
		},
		{
			name: "Successfully case",
			fields: fields{
				metrics: map[string]models.Metric{
					"metric1": {
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"metric2": {
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				metrics: models.MetricsList{
					{
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					{
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{
				metrics: test.fields.metrics,
			}
			got, err := s.GetAllMetrics(test.args.ctx)
			assert.Equal(t, &test.want.metrics, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_UpdateMetrics(t *testing.T) {
	dOld := int64(50)
	dNew := int64(100)
	vOld := float64(50.05)
	vNew := float64(100.1)
	type fields struct {
		metrics  map[string]models.Metric
		filePath string
		syncMode bool
	}
	type args struct {
		ctx     context.Context
		metrics models.MetricsList
	}
	type want struct {
		metrics map[string]models.Metric
		err     error
	}
	tests := []struct {
		fields fields
		name   string
		want   want
		args   args
	}{
		{
			name: "Update empty store without sync",
			fields: fields{
				metrics:  map[string]models.Metric{},
				syncMode: false,
				filePath: fmt.Sprintf("/tmp/%d", time.Now().UnixMicro()),
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					{
						Value: &vNew,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					{
						Value: nil,
						Delta: &dNew,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			want: want{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: &vNew,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"bab0d7e52057f00e3cf6b851bf72f4c1": {
						Value: nil,
						Delta: &dNew,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				err: nil,
			},
		},
		{
			name: "Update empty list without sync",
			fields: fields{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: &vOld,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"bab0d7e52057f00e3cf6b851bf72f4c1": {
						Value: nil,
						Delta: &dOld,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				syncMode: false,
				filePath: fmt.Sprintf("/tmp/%d", time.Now().UnixMicro()),
			},
			args: args{
				ctx:     context.Background(),
				metrics: models.MetricsList{},
			},
			want: want{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: &vOld,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"bab0d7e52057f00e3cf6b851bf72f4c1": {
						Value: nil,
						Delta: &dOld,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				err: nil,
			},
		},
		{
			name: "Successfully case with sync mode",
			fields: fields{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: &vOld,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"bab0d7e52057f00e3cf6b851bf72f4c1": {
						Value: nil,
						Delta: &dOld,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				syncMode: false,
				filePath: fmt.Sprintf("/tmp/%d", time.Now().UnixMicro()),
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					{
						Value: &vNew,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					{
						Value: nil,
						Delta: &dNew,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			want: want{
				metrics: map[string]models.Metric{
					"c917499366554154b438bb8259183adc": {
						Value: &vNew,
						Delta: nil,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"bab0d7e52057f00e3cf6b851bf72f4c1": {
						Value: nil,
						Delta: &dNew,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{
				metrics:  test.fields.metrics,
				syncMode: test.fields.syncMode,
				filePath: test.fields.filePath,
			}
			err := s.UpdateMetrics(test.args.ctx, test.args.metrics)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.metrics, s.metrics)
		})
	}
}

func TestStore_Save(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		metrics  map[string]models.Metric
		filePath string
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err       error
		content   []byte
		fileExist bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "successfully case",
			fields: fields{
				metrics: map[string]models.Metric{
					"abc123": {
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 1",
					},
					"abc124": {
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 2",
					},
				},
				filePath: "/tmp/TestStore_Save.bac",
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:       nil,
				fileExist: true,
				content:   []byte("[{\"value\":100.1,\"type\":\"gauge\",\"id\":\"metric 1\"},{\"delta\":100,\"type\":\"counter\",\"id\":\"metric 2\"}]"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Remove(test.fields.filePath)
			s := &Store{
				metrics:  test.fields.metrics,
				filePath: test.fields.filePath,
			}
			err := s.Save(test.args.ctx)
			require.Equal(t, test.want.err, err)
			if test.want.fileExist {
				require.FileExists(t, test.fields.filePath)
				c, err := os.ReadFile(test.fields.filePath)
				require.Nil(t, err)
				assert.Equal(t, test.want.content, c)
			}
		})
	}
}
