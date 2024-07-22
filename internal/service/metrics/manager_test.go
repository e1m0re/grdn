package metrics

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/storage/store"
	"github.com/e1m0re/grdn/internal/storage/store/mocks"
)

func Test_metricsManager_GetAllMetrics(t *testing.T) {
	delta := int64(100)
	value := float64(100.10)
	type args struct {
		ctx context.Context
	}
	type want struct {
		result *models.MetricsList
		errMsg string
	}
	tests := []struct {
		name      string
		args      args
		mockStore func() store.Store
		want      want
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetAllMetrics", mock.Anything).
					Return(nil, errors.New("something wrong"))
				return mockStore
			},
			want: want{
				result: nil,
				errMsg: "something wrong",
			},
		},
		{
			name: "Empty result",
			args: args{
				ctx: context.Background(),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetAllMetrics", mock.Anything).
					Return(&models.MetricsList{}, nil)
				return mockStore
			},
			want: want{
				result: &models.MetricsList{},
				errMsg: "",
			},
		},
		{
			name: "Empty result",
			args: args{
				ctx: context.Background(),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetAllMetrics", mock.Anything).
					Return(&models.MetricsList{
						&models.Metric{
							Value: nil,
							Delta: &delta,
							MType: models.CounterType,
							ID:    "metric1",
						},
						&models.Metric{
							Value: &value,
							Delta: nil,
							MType: models.GaugeType,
							ID:    "metric2",
						},
					}, nil)
				return mockStore
			},
			want: want{
				result: &models.MetricsList{
					&models.Metric{
						Value: nil,
						Delta: &delta,
						MType: models.CounterType,
						ID:    "metric1",
					},
					&models.Metric{
						Value: &value,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric2",
					},
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm := NewMetricsManager(test.mockStore())
			got, err := mm.GetAllMetrics(test.args.ctx)
			if len(test.want.errMsg) > 0 {
				require.Errorf(t, err, test.want.errMsg)
			}
			require.Equal(t, test.want.result, got)
		})
	}
}

func Test_metricsManager_GetMetric(t *testing.T) {
	delta := int64(100)
	value := float64(100.10)
	type args struct {
		ctx  context.Context
		Type models.MetricType
		Name models.MetricName
	}
	type want struct {
		result *models.Metric
		errMsg string
	}
	tests := []struct {
		name      string
		args      args
		mockStore func() store.Store
		want      want
	}{
		{
			name: "something wrong",
			args: args{
				ctx:  context.Background(),
				Type: models.GaugeType,
				Name: "metric1",
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, errors.New("something wrong"))

				return mockStore
			},
			want: want{
				result: nil,
				errMsg: "something wrong",
			},
		},
		{
			name: "metric not found",
			args: args{
				ctx:  context.Background(),
				Type: models.GaugeType,
				Name: "metric1",
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, nil)

				return mockStore
			},
			want: want{
				result: nil,
				errMsg: "",
			},
		},
		{
			name: "returns gauge metric",
			args: args{
				ctx:  context.Background(),
				Type: models.GaugeType,
				Name: "metric1",
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						Value: &value,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric1",
					}, nil)

				return mockStore
			},
			want: want{
				result: &models.Metric{
					Value: &value,
					Delta: nil,
					MType: models.GaugeType,
					ID:    "metric1",
				},
				errMsg: "",
			},
		},
		{
			name: "returns counter metric",
			args: args{
				ctx:  context.Background(),
				Type: models.CounterType,
				Name: "metric1",
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						Value: nil,
						Delta: &delta,
						MType: models.CounterType,
						ID:    "metric1",
					}, nil)

				return mockStore
			},
			want: want{
				result: &models.Metric{
					Value: nil,
					Delta: &delta,
					MType: models.CounterType,
					ID:    "metric1",
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm := NewMetricsManager(test.mockStore())
			got, err := mm.GetMetric(test.args.ctx, test.args.Type, test.args.Name)
			if len(test.want.errMsg) > 0 {
				require.Errorf(t, err, test.want.errMsg)
			}
			require.Equal(t, test.want.result, got)
		})
	}
}

func Test_metricsManager_UpdateMetric(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		mockStore func() store.Store
	}
	type args struct {
		ctx    context.Context
		metric models.Metric
	}
	type want struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "store.GetMetric failed",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, errors.New("something wrong"))
					return mockStore
				},
			},
			args: args{
				ctx:    context.Background(),
				metric: models.Metric{},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "unknown metric type",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil)

					return mockStore
				},
			},
			args: args{
				ctx:    context.Background(),
				metric: models.Metric{},
			},
			want: want{
				err: errors.New("unknown metric type"),
			},
		},
		{
			name: "store.UpdateMetrics failed",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil)
					mockStore.
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(errors.New("something wrong"))

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metric: models.Metric{
					MType: models.GaugeType,
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "Update gauge metric (successfully case)",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil).
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(nil)

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metric: models.Metric{
					MType: models.GaugeType,
					Value: &v,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Update new counter metric (successfully case)",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(&models.Metric{
							Value: nil,
							Delta: &d,
							MType: models.CounterType,
							ID:    "metric 1",
						}, nil).
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(nil)

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metric: models.Metric{
					ID:    "metric 1",
					MType: models.CounterType,
					Delta: &d,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Update counter metric (successfully case)",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil).
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(nil)

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metric: models.Metric{
					MType: models.CounterType,
					Delta: &d,
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm := &metricsManager{
				store: test.fields.mockStore(),
			}
			err := mm.UpdateMetric(test.args.ctx, test.args.metric)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_metricsManager_UpdateMetrics(t *testing.T) {
	d := int64(100)
	v := float64(100.1)
	type fields struct {
		mockStore func() store.Store
	}
	type args struct {
		ctx     context.Context
		metrics models.MetricsList
	}
	type want struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Successfully update (empty args list)",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)

					return mockStore
				},
			},
			args: args{
				ctx:     context.Background(),
				metrics: make(models.MetricsList, 0),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Successfully update (empty args list)",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil).
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(nil)

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					&models.Metric{
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "GetMetric failed",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, errors.New("something wrong"))

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					&models.Metric{
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "UpdateMetrics failed",
			fields: fields{
				mockStore: func() store.Store {
					mockStore := mocks.NewStore(t)
					mockStore.
						On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(nil, nil).
						On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
						Return(errors.New("something wrong"))

					return mockStore
				},
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						Value: nil,
						Delta: &d,
						MType: models.CounterType,
						ID:    "metric 1",
					},
					&models.Metric{
						Value: &v,
						Delta: nil,
						MType: models.GaugeType,
						ID:    "metric 2",
					},
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm := &metricsManager{
				store: test.fields.mockStore(),
			}
			err := mm.UpdateMetrics(test.args.ctx, test.args.metrics)
			assert.Equal(t, test.want.err, err)
		})
	}
}
