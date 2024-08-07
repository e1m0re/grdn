package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/service"
	"github.com/e1m0re/grdn/internal/service/metrics/mocks"
)

func TestHandler_getMetricValueV2(t *testing.T) {
	delta := int64(100)
	value := 100.123456
	type args struct {
		body   string
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
		mockServices func() *service.ServerServices
		args         args
		want         want
	}{
		{
			name: "Invalid method",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
		{
			name: "Invalid Body",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "error parsing body: unexpected end of JSON input\n",
			},
		},
		{
			name: "Unknown metric",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, nil)

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				body:   "{\"id\":\"metricId\",\"type\":\"metricType\"}",
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusNotFound,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "Not found.\n",
			},
		},
		{
			name: "GetMetric failed",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, errors.New("something wrong"))

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				body:   "{\"id\":\"metricId\",\"type\":\"metricType\"}",
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "something wrong\n",
			},
		},
		{
			name: "Successfully test (Counter metric)",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						ID:    "metricId",
						MType: models.CounterType,
						Delta: &delta,
						Value: nil,
					}, nil)

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				body:   fmt.Sprintf("{\"id\":\"metricId\",\"type\":\"%s\",\"delta\":%d}", models.CounterType, delta),
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedHeaders:      map[string]string{"Content-Type": "application/json"},
				expectedResponseBody: fmt.Sprintf("{\"delta\":%d,\"type\":\"%s\",\"id\":\"metricId\"}", delta, models.CounterType),
			},
		},
		{
			name: "Successfully test (GaugeType metric)",
			mockServices: func() *service.ServerServices {
				mockMetricsManager := mocks.NewManager(t)
				mockMetricsManager.
					On("GetMetric", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(&models.Metric{
						ID:    "metricId",
						MType: models.GaugeType,
						Delta: nil,
						Value: &value,
					}, nil)

				return &service.ServerServices{
					MetricsManager: mockMetricsManager,
				}
			},
			args: args{
				ctx:    context.Background(),
				body:   fmt.Sprintf("{\"id\":\"metricId\",\"type\":\"%s\",\"value\":%f}", models.GaugeType, value),
				method: http.MethodPost,
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedHeaders:      map[string]string{"Content-Type": "application/json"},
				expectedResponseBody: fmt.Sprintf("{\"value\":%f,\"type\":\"%s\",\"id\":\"metricId\"}", value, models.GaugeType),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("", "")

			req, err := http.NewRequestWithContext(test.args.ctx, test.args.method, "/value", bytes.NewReader([]byte(test.args.body)))
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
