package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateMetricHandler(t *testing.T) {
	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name: "test invalid path",
			args: args{
				method: http.MethodPost,
				path:   "/",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid method",
			args: args{
				method: http.MethodDelete,
				path:   "/update",
			},
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name: "test invalid metric type",
			args: args{
				method: http.MethodPost,
				path:   "/update/unknown",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid metric name",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/unknown",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test empty metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/Alloc",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid gauge metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/Alloc/asd",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid counter metric value",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter/PollCount/asd",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test success update gauge metric",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/Alloc/123.12",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "test success update counter metric",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter/PollCount/123",
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			response := httptest.NewRecorder()
			updateMetricHandler(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)
			fmt.Printf("%v", response)
		})
	}
}
