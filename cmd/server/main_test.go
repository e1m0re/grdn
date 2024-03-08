package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e1m0re/grdn/internal/storage"

	"github.com/stretchr/testify/assert"
)

func Test_isValidMetricName(t *testing.T) {
	type args struct {
		mType storage.MetricsType
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check gauge type",
			args: args{
				mType: storage.GuageType,
				value: storage.Alloc,
			},
			want: true,
		},
		{
			name: "test check counter type",
			args: args{
				mType: storage.CounterType,
				value: "PollCount",
			},
			want: true,
		},
		{
			name: "test empty type",
			args: args{
				mType: "",
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown type",
			args: args{
				mType: "Unknown",
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown value for guage type",
			args: args{
				mType: storage.GuageType,
				value: "PollCount",
			},
			want: false,
		},
		{
			name: "test unknown value for Counter type",
			args: args{
				mType: storage.CounterType,
				value: "PollCount1",
			},
			want: false,
		},
		{
			name: "test empty type",
			args: args{
				mType: storage.CounterType,
				value: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isValidMetricName(tt.args.mType, tt.args.value))
		})
	}
}

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
				path:   "/upload",
			},
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name: "test invalid metric type",
			args: args{
				method: http.MethodPost,
				path:   "/upload/unknown",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid metric name",
			args: args{
				method: http.MethodPost,
				path:   "/upload/guage/unknown",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "test empty metric value",
			args: args{
				method: http.MethodPost,
				path:   "/upload/guage/Alloc",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid guage metric value",
			args: args{
				method: http.MethodPost,
				path:   "/upload/guage/Alloc/asd",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test invalid counter metric value",
			args: args{
				method: http.MethodPost,
				path:   "/upload/counter/PollCount/asd",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test success update guage metric",
			args: args{
				method: http.MethodPost,
				path:   "/upload/guage/Alloc/123.12",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "test success update counter metric",
			args: args{
				method: http.MethodPost,
				path:   "/upload/counter/PollCount/123",
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
		})
	}
}
