package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateMetricHandler(t *testing.T) {
	type args struct {
		method       string
		path         string
		routerParams map[string]string
	}

	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name: "test invalid path",
			args: args{
				method:       http.MethodPost,
				path:         "/",
				routerParams: make(map[string]string),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid method",
			args: args{
				method:       http.MethodDelete,
				path:         "/update",
				routerParams: make(map[string]string),
			},
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name: "test invalid metric type",
			args: args{
				method: http.MethodPost,
				path:   "/update/unknown",
				routerParams: map[string]string{
					"mType": "unknown",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test empty metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/Alloc",
				routerParams: map[string]string{
					"mType": "gauge",
					"mName": "Alloc",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid gauge metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/Alloc/asd",
				routerParams: map[string]string{
					"mType":  "gauge",
					"mName":  "Alloc",
					"mValue": "asd",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid counter metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter/PollCount/asd",
				routerParams: map[string]string{
					"mType":  "counter",
					"mName":  "PollCount",
					"mValue": "asd",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test success update gauge metric",
			args: args{
				method: http.MethodPost,
				path:   "http://localhost/update/gauge/Alloc/123.12",
				routerParams: map[string]string{
					"mType":  "gauge",
					"mName":  "Alloc",
					"mValue": "123.12",
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "test success update counter metric",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter/PollCount/123",
				routerParams: map[string]string{
					"mType":  "counter",
					"mName":  "PollCount",
					"mValue": "123",
				},
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			response := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			for key, value := range tt.args.routerParams {
				rctx.URLParams.Add(key, value)
			}

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			updateMetricHandler(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
