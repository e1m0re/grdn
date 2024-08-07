package memory

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
)

var (
	delta = int64(100)
	value = float64(100.1)
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
						Delta: &delta,
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
					Delta: &delta,
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
	type fields struct {
		metrics map[string]models.Metric
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err     error
		metrics map[string]models.Metric
	}
	tests := []struct {
		want   want
		fields fields
		args   args
		name   string
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
				metrics: make(map[string]models.Metric),
				err:     nil,
			},
		},
		{
			name: "Successfully case",
			fields: fields{
				metrics: map[string]models.Metric{
					"metric1": {
						Value: nil,
						Delta: &delta,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"metric2": {
						Value: &value,
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
				metrics: map[string]models.Metric{
					"metric 1": {
						Value: nil,
						Delta: &delta,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					"metric 2": {
						Value: &value,
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
			m := make(map[string]models.Metric)
			for _, metric := range *got {
				m[metric.ID] = *metric
			}
			assert.Equal(t, test.want.metrics, m)
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
					"9e646d6d5855fbdadcaab202f1748505": {
						Value: &value,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 1",
					},
					"3e0673a56ff12916a6293fa5a1bfc2db": {
						Value: nil,
						Delta: &delta,
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

func TestStore_Restore(t *testing.T) {
	type fields struct {
		metrics  map[string]models.Metric
		filePath string
		content  []byte
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		metrics map[string]models.Metric
		err     error
	}
	tests := []struct {
		name   string
		args   args
		want   want
		fields fields
	}{
		{
			name: "successfully case",
			fields: fields{
				metrics:  make(map[string]models.Metric),
				filePath: "/tmp/TestStore_Restore.bac",
				content:  []byte("[{\"delta\":100,\"type\":\"counter\",\"id\":\"metric 2\"},{\"value\":100.1,\"type\":\"gauge\",\"id\":\"metric 1\"}]"),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
				metrics: map[string]models.Metric{
					"9e646d6d5855fbdadcaab202f1748505": {
						Value: &value,
						MType: models.GaugeType,
						ID:    "metric 1",
					},
					"3e0673a56ff12916a6293fa5a1bfc2db": {
						Delta: &delta,
						MType: models.CounterType,
						ID:    "metric 2",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Create(test.fields.filePath)
			require.Nil(t, err)
			_, err = f.Write(test.fields.content)
			require.Nil(t, err)
			s := &Store{
				metrics:  test.fields.metrics,
				filePath: test.fields.filePath,
			}
			err = s.Restore(test.args.ctx)
			require.Equal(t, test.want.err, err)
			if err == nil {
				require.Equal(t, len(test.want.metrics), len(s.metrics))
				for key, metric := range test.want.metrics {
					m, ok := s.metrics[key]
					require.True(t, ok)
					require.Equal(t, metric, m)
				}
			}
		})
	}
}

func TestNewStore(t *testing.T) {
	type args struct {
		ctx      context.Context
		filePath string
		syncMode bool
	}
	type want struct {
		str *Store
		err error
	}
	tests := []struct {
		want want
		name string
		args args
	}{
		{
			name: "Successfully case without restore",
			args: args{
				ctx:      context.Background(),
				filePath: "",
				syncMode: false,
			},
			want: want{
				str: &Store{
					metrics:  make(map[string]models.Metric),
					filePath: "",
					syncMode: false,
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewStore(test.args.ctx, test.args.filePath, test.args.syncMode)
			require.Equal(t, test.want.err, err)
			//assert.Implements(t, (*store.Store)(nil), got)
			assert.Equal(t, test.want.str.metrics, got.metrics)
			assert.Equal(t, test.want.str.filePath, got.filePath)
			assert.Equal(t, test.want.str.syncMode, got.syncMode)
		})
	}
}
