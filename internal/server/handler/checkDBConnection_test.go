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

func TestHandler_checkDBConnection(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
	}
	type want struct {
		expectedStatusCode   int
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
				expectedResponseBody: "",
			},
		},
		{
			name: "Check connection failed with error",
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			mockServices: func() *service.Services {
				mockMetricsService := mockservice.NewMetricsService(t)
				mockMetricsService.
					On("PingDB", mock.Anything).
					Return(fmt.Errorf("something wrong"))

				return &service.Services{
					MetricsService: mockMetricsService,
				}
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name: "Successfully test",
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
			},
			mockServices: func() *service.Services {
				mockMetricsService := mockservice.NewMetricsService(t)
				mockMetricsService.
					On("PingDB", mock.Anything).
					Return(nil)

				return &service.Services{
					MetricsService: mockMetricsService,
				}
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter("")

			req, err := http.NewRequestWithContext(test.args.ctx, test.args.method, "/ping", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
