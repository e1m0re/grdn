package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/service"
	mockservice "github.com/e1m0re/grdn/internal/service/mocks"
)

func TestHandler_updateMetricsList(t *testing.T) {
	type args struct {
		body   string
		ctx    context.Context
		method string
		path   string
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
				mockMetricService := mockservice.NewMetricsService(t)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/updates",
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
		{
			name: "Invalid Body",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				body:   "",
				ctx:    context.Background(),
				method: http.MethodPost,
				path:   "/updates",
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "unexpected end of JSON input\n",
			},
		},
		{
			name: "UpdateMetrics failed",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
					Return(fmt.Errorf("something wrong"))

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				body:   "[{\"id\":\"metricId\",\"type\":\"metricType\"}]",
				ctx:    context.Background(),
				method: http.MethodPost,
				path:   "/updates",
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
		{
			name: "Successfully test",
			mockServices: func() *service.Services {
				mockMetricService := mockservice.NewMetricsService(t)
				mockMetricService.
					On("UpdateMetrics", mock.Anything, mock.AnythingOfType("models.MetricsList")).
					Return(nil)

				return &service.Services{
					MetricsService: mockMetricService,
				}
			},
			args: args{
				body:   "[{\"id\":\"metricId\",\"type\":\"metricType\"}]",
				ctx:    context.Background(),
				method: http.MethodPost,
				path:   "/updates",
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedHeaders:      make(map[string]string),
				expectedResponseBody: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("")

			req, err := http.NewRequestWithContext(test.args.ctx, test.args.method, test.args.path, bytes.NewReader([]byte(test.args.body)))
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
