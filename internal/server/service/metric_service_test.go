package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage"
	"github.com/e1m0re/grdn/internal/server/storage/mocks"
)

func Test_metricService_GetMetric(t *testing.T) {
	delta := int64(100)
	value := float64(100.11)

	type args struct {
		ctx   context.Context
		mType models.MetricType
		mName models.MetricName
	}
	type want struct {
		result *models.Metric
		errMsg string
	}
	tests := []struct {
		name     string
		getMocks func() storage.Storage
		args     args
		want     want
	}{
		{
			name: "ErrUnknownMetric error",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, storage.ErrUnknownMetric)

				return mockStorage
			},
			want: want{
				result: nil,
				errMsg: storage.ErrUnknownMetric.Error(),
			},
		},
		{
			name: "get Counter result",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						Value: nil,
						Delta: &delta,
						MType: models.CounterType,
						ID:    models.PollCount,
					}, nil)

				return mockStorage
			},
			want: want{
				result: &models.Metric{
					Value: nil,
					Delta: &delta,
					MType: models.CounterType,
					ID:    models.PollCount,
				},
				errMsg: "",
			},
		},
		{
			name: "get Gauge result",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						Value: &value,
						Delta: nil,
						MType: models.GaugeType,
						ID:    models.Alloc,
					}, nil)

				return mockStorage
			},
			want: want{
				result: &models.Metric{
					Value: &value,
					Delta: nil,
					MType: models.GaugeType,
					ID:    models.Alloc,
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := &metricService{
				store: test.getMocks(),
			}

			gotMetric, err := ms.GetMetric(test.args.ctx, test.args.mType, test.args.mName)
			assert.Equal(t, test.want.result, gotMetric)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, err, test.want.errMsg)
			}
		})
	}
}

func Test_metricService_GetMetricsList(t *testing.T) {
	delta := int64(100)
	value := float64(100.11)

	type args struct {
		ctx context.Context
	}
	type want struct {
		result []string
		errMsg string
	}
	tests := []struct {
		name     string
		getMocks func() storage.Storage
		args     args
		want     want
	}{
		{
			name: "something wrong",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetricsList", mock.Anything).
					Return(nil, fmt.Errorf("something wrong"))

				return mockStorage
			},
			want: want{
				result: nil,
				errMsg: "something wrong",
			},
		},
		{
			name: "empty list",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetricsList", mock.Anything).
					Return(make([]string, 0), nil)

				return mockStorage
			},
			want: want{
				result: make([]string, 0),
				errMsg: "",
			},
		},
		{
			name: "success case",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("GetMetricsList", mock.Anything).
					Return([]string{
						fmt.Sprintf("TotalAlloc: %f", value),
						fmt.Sprintf("PollCount: %d", delta),
					}, nil)

				return mockStorage
			},
			want: want{
				result: []string{
					fmt.Sprintf("TotalAlloc: %f", value),
					fmt.Sprintf("PollCount: %d", delta),
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := &metricService{
				store: test.getMocks(),
			}
			got, err := ms.GetMetricsList(test.args.ctx)
			assert.Equal(t, test.want.result, got)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, err, test.want.errMsg)
			}
		})
	}
}

func Test_metricService_UpdateMetric(t *testing.T) {
	type args struct {
		ctx    context.Context
		metric models.Metric
	}
	type want struct {
		errMsg string
	}
	tests := []struct {
		name     string
		getMocks func() storage.Storage
		args     args
		want     want
	}{
		{
			name: "something went wrong",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("UpdateMetric", mock.Anything, mock.AnythingOfType("models.Metric")).
					Return(fmt.Errorf("something wrong"))

				return mockStorage
			},
			args: args{
				ctx:    context.Background(),
				metric: models.Metric{},
			},
			want: want{
				errMsg: "something wrong",
			},
		},
		{
			name: "successfully case",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("UpdateMetric", mock.Anything, mock.AnythingOfType("models.Metric")).
					Return(nil)

				return mockStorage
			},
			args: args{
				ctx:    context.Background(),
				metric: models.Metric{},
			},
			want: want{
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := &metricService{
				store: test.getMocks(),
			}

			err := ms.UpdateMetric(test.args.ctx, test.args.metric)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, err, test.want.errMsg)
			}
		})
	}
}

func Test_metricService_UpdateMetrics(t *testing.T) {
	type args struct {
		ctx     context.Context
		metrics models.MetricsList
	}
	type want struct {
		errMsg string
	}
	tests := []struct {
		name     string
		getMocks func() storage.Storage
		args     args
		want     want
	}{
		{
			name: "something went wrong",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
					Return(fmt.Errorf("something wrong"))

				return mockStorage
			},
			args: args{
				ctx:     context.Background(),
				metrics: make(models.MetricsList, 0),
			},
			want: want{
				errMsg: "something wrong",
			},
		},
		{
			name: "successfully case",
			getMocks: func() storage.Storage {
				mockStorage := mocks.NewStorage(t)
				mockStorage.
					On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
					Return(nil)

				return mockStorage
			},
			args: args{
				ctx:     context.Background(),
				metrics: make(models.MetricsList, 0),
			},
			want: want{
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := &metricService{
				store: test.getMocks(),
			}

			err := ms.UpdateMetrics(test.args.ctx, test.args.metrics)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, err, test.want.errMsg)
			}
		})
	}
}
