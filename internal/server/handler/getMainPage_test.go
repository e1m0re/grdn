package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/service"
	"github.com/e1m0re/grdn/internal/server/service/metrics/mocks"
)

func TestHandler_getMainPage(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
	}
	type want struct {
		expectedHeaders      map[string]string
		expectedResponseBody string
		expectedStatusCode   int
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
				mockMetricsManager := mocks.NewManager(t)

				return &service.Services{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name: "Request failed",
			mockServices: func() *service.Services {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetAllMetrics", mock.Anything).
					Return(nil, errors.New("something wrong"))

				return &service.Services{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
		{
			name: "Empty metrics list",
			mockServices: func() *service.Services {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetAllMetrics", mock.Anything).
					Return(&models.MetricsList{}, nil)

				return &service.Services{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedHeaders:      map[string]string{"Content-Type": "text/html"},
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "",
			},
		},
		{
			name: "Successful test",
			mockServices: func() *service.Services {
				value := float64(100.100)
				metric1 := &models.Metric{
					Value: &value,
					Delta: nil,
					MType: models.GaugeType,
					ID:    "metric1",
				}
				delta := int64(100)
				metric2 := &models.Metric{
					Value: nil,
					Delta: &delta,
					MType: models.CounterType,
					ID:    "metric2",
				}
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetAllMetrics", mock.Anything).
					Return(&models.MetricsList{metric1, metric2}, nil)

				return &service.Services{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedHeaders:      map[string]string{"Content-Type": "text/html"},
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "metric1: 100.1\r\nmetric2: 100\r\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("", "")

			req, err := http.NewRequestWithContext(test.args.ctx, test.args.method, "/", nil)
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
