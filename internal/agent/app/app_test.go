package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/service"
	mocks3 "github.com/e1m0re/grdn/internal/service/apiclient/mocks"
	mocks2 "github.com/e1m0re/grdn/internal/service/encryption/mocks"
	"github.com/e1m0re/grdn/internal/service/monitor/mocks"
)

func TestApp_updateDataWorker(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err     error
		metrics models.MetricsList
	}
	tests := []struct {
		args    args
		mockApp func() *app
		name    string
		want    want
	}{
		{
			name: "successfully case",
			mockApp: func() *app {
				mockAPIClient := mocks3.NewAPIClient(t)
				mockAPIClient.On("SendMetricsData", &[]byte{}).Return(nil).Maybe()

				mockMonitor := mocks.NewMonitor(t)
				mockMonitor.On("UpdateData", mock.Anything).Return(nil)
				mockMonitor.On("UpdateGOPS", mock.Anything).Return(nil)
				mockMonitor.On("GetMetricsList").Return(make(models.MetricsList, 0))

				mockEncryptor := mocks2.NewEncryptor(t)
				mockEncryptor.
					On("Encrypt", []byte{0x5b, 0x5d}).
					Return(make([]byte, 0), nil)

				return &app{
					apiClient: mockAPIClient,
					monitor:   mockMonitor,
					encryptor: mockEncryptor,
					cfg: &config.Config{
						PollInterval: time.Second,
						RateLimit:    3,
					},
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:     nil,
				metrics: make(models.MetricsList, 0),
			},
		},
		{
			name: "UpdateGOPS failed",
			mockApp: func() *app {
				mockAPIClient := mocks3.NewAPIClient(t)
				mockAPIClient.On("SendMetricsData", &[]byte{}).Return(nil).Maybe()

				mockMonitor := mocks.NewMonitor(t)
				mockMonitor.On("UpdateData", mock.Anything).Return(nil)
				mockMonitor.On("UpdateGOPS", mock.Anything).Return(errors.New("something wrong"))
				mockMonitor.On("GetMetricsList").Return(make(models.MetricsList, 0))

				mockEncryptor := mocks2.NewEncryptor(t)
				mockEncryptor.
					On("Encrypt", []byte{0x5b, 0x5d}).
					Return(make([]byte, 0), nil)

				return &app{
					apiClient: mockAPIClient,
					monitor:   mockMonitor,
					encryptor: mockEncryptor,
					cfg:       &config.Config{PollInterval: time.Second},
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:     nil,
				metrics: make(models.MetricsList, 0),
			},
		},
		{
			name: "encryption failed",
			mockApp: func() *app {
				mockAPIClient := mocks3.NewAPIClient(t)
				mockAPIClient.On("SendMetricsData", &[]byte{}).Return(nil).Maybe()

				mockMonitor := mocks.NewMonitor(t)
				mockMonitor.On("UpdateData", mock.Anything).Return(nil)
				mockMonitor.On("UpdateGOPS", mock.Anything).Return(nil)
				mockMonitor.On("GetMetricsList").Return(make(models.MetricsList, 0))

				mockEncryptor := mocks2.NewEncryptor(t)
				mockEncryptor.
					On("Encrypt", []byte{0x5b, 0x5d}).
					Return(nil, errors.New("something wrong"))

				return &app{
					apiClient: mockAPIClient,
					monitor:   mockMonitor,
					encryptor: mockEncryptor,
					cfg:       &config.Config{PollInterval: time.Second},
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:     nil,
				metrics: make(models.MetricsList, 0),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := test.mockApp()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				err := app.Start(ctx)
				if err != nil {
					panic(err)
				}
			}()

			<-time.After(time.Second * 3)
			metrics := app.monitor.GetMetricsList()
			assert.Equal(t, test.want.metrics, metrics)
		})
	}
}

func TestNewApp(t *testing.T) {
	cfg := &config.Config{}
	services, err := service.NewAgentServices(cfg)
	require.Nil(t, err)
	app := NewApp(cfg, services)
	assert.Implements(t, (*App)(nil), app)
}
