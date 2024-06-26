package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/service"
	mockservice "github.com/e1m0re/grdn/internal/service/mocks"
	"github.com/e1m0re/grdn/internal/storage"
)

func TestHandler_getMetricValue(t *testing.T) {
	delta := int64(100)
	value := 100.123456
	type args struct {
		ctx    context.Context
		method string
	}
	type want struct {
		expectedStatusCode   int
		expectedHeaders      map[string]string
		expectedResponseBody string
	}
	tests := []struct {
		name         string
		mockServices func() *service.Services
		args         args
		want         want
	}{
		{
			name: "Invalid method",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
		{
			name: "Unknown metric",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, storage.ErrUnknownMetric)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusNotFound,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "Not found.\n",
			},
		},
		{
			name: "GetMetric failed",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, fmt.Errorf("something wrong"))

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "something wrong\n",
			},
		},
		{
			name: "Successfully test (Counter metric)",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						ID:    "metricId",
						MType: models.CounterType,
						Delta: &delta,
						Value: nil,
					}, nil)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedHeaders:      map[string]string{"Content-Type": "text/html"},
				expectedResponseBody: fmt.Sprintf("%d", delta),
			},
		},
		{
			name: "Successfully test (Gauge metric)",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						ID:    "metricId",
						MType: models.GaugeType,
						Delta: nil,
						Value: &value,
					}, nil)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedHeaders:      map[string]string{"Content-Type": "text/html"},
				expectedResponseBody: fmt.Sprintf("%f", value),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("")

			req, err := http.NewRequestWithContext(test.args.ctx, test.args.method, "/value/mType/mName", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			for key, value := range test.want.expectedHeaders {
				require.Equal(t, value, rr.Header().Get(key))
			}
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
