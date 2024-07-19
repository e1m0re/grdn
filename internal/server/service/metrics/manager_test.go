package metrics

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage/store"
	"github.com/e1m0re/grdn/internal/server/storage/store/mocks"
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
					Return(nil, fmt.Errorf("something wrong"))
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
					Return(nil, fmt.Errorf("something wrong"))

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
