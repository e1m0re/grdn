package handler

import (
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
				expectedResponseBody: "",
			},
		},
		{
			name: "Request failed",
			mockServices: func() *service.Services {
				mockMetricsService := mockservice.NewMetricsService(t)
				mockMetricsService.
					On("GetMetricsList", mock.Anything).
					Return(make([]string, 0), fmt.Errorf("something wrong"))

				return &service.Services{
					MetricsService: mockMetricsService,
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
				mockMetricsService := mockservice.NewMetricsService(t)
				mockMetricsService.
					On("GetMetricsList", mock.Anything).
					Return(make([]string, 0), nil)

				return &service.Services{
					MetricsService: mockMetricsService,
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
			name: "Successfult test",
			mockServices: func() *service.Services {
				mockMetricsService := mockservice.NewMetricsService(t)
				mockMetricsService.
					On("GetMetricsList", mock.Anything).
					Return([]string{"metric1", "metric2", "metric3"}, nil)

				return &service.Services{
					MetricsService: mockMetricsService,
				}
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			want: want{
				expectedHeaders:      map[string]string{"Content-Type": "text/html"},
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "metric1\r\nmetric2\r\nmetric3\r\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("")

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
