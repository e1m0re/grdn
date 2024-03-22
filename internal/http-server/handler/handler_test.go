package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e1m0re/grdn/internal/storage"
)

func TestHandler_updateMetricHandler(t *testing.T) {
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

	store := storage.NewMemStorage()
	handler := http.HandlerFunc(NewHandler(store).UpdateMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			response := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			for key, value := range tt.args.routerParams {
				rctx.URLParams.Add(key, value)
			}

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			handler.ServeHTTP(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func TestHandler_UpdateMetrics(t *testing.T) {
	type fields struct {
		store *storage.MemStorage
	}

	type args struct {
		request *http.Request
	}

	type want struct {
		statusCode int
		store      *storage.MemStorage
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "test invalid method",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/update", nil),
			},
			want: want{statusCode: http.StatusMethodNotAllowed},
		},
		{
			name:   "test empty body",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", nil),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test invalid body",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{123:123}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body without ID",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"type":"counter","delta":10}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body invalid type",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"counter1","delta":10}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body invalid gauge value 1",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"gauge"}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body invalid gauge value 2",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"gauge","value":"10"}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body invalid counter value 1",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"counter"}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test body invalid counter value 2",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"counter","delta":"10"}`)),
			},
			want: want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "test success update gauge metric",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"gauge","value":10.10}`)),
			},
			want: want{
				statusCode: http.StatusOK,
				store: &storage.MemStorage{
					Gauges: map[storage.GaugeName]storage.GaugeDateType{
						"metric": 10.10,
					},
					Counters: make(map[storage.CounterName]storage.CounterDateType),
				},
			},
		},
		{
			name:   "test success update counter metric",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"metric","type":"counter","delta":10}`)),
			},
			want: want{
				statusCode: http.StatusOK,
				store: &storage.MemStorage{
					Gauges: make(map[storage.GaugeName]storage.GaugeDateType),
					Counters: map[storage.CounterName]storage.CounterDateType{
						"metric": 10,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := http.HandlerFunc(NewHandler(tt.fields.store).UpdateMetrics)

			response := httptest.NewRecorder()

			handler.ServeHTTP(response, tt.args.request)
			require.Equal(t, tt.want.statusCode, response.Code)

			if tt.want.statusCode == http.StatusOK {
				assert.Equal(t, tt.want.store, tt.fields.store)
			}
		})
	}
}

func TestHandler_GetMetricValue(t *testing.T) {

	type fields struct {
		store *storage.MemStorage
	}

	type args struct {
		request      *http.Request
		routerParams map[string]string
	}

	type want struct {
		statusCode int
		content    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "test invalid method",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request:      httptest.NewRequest(http.MethodPost, "/value", nil),
				routerParams: nil,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "test invalid metric type",
			fields: fields{store: storage.NewMemStorage()},
			args: args{
				request:      httptest.NewRequest(http.MethodGet, "/value/something", nil),
				routerParams: nil,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "test unknown metric",
			fields: fields{
				store: &storage.MemStorage{
					Gauges: make(map[storage.GaugeName]storage.GaugeDateType),
					Counters: map[storage.CounterName]storage.CounterDateType{
						"counter1": 1984,
					},
				},
			},
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/value/counter/counter2", nil),
				routerParams: map[string]string{
					"mType": "counter",
					"mName": "counter2",
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "test of successfully getting gauge metric",
			fields: fields{
				store: &storage.MemStorage{
					Gauges: map[storage.GaugeName]storage.GaugeDateType{
						"Alloc": 10.10,
					},
					Counters: make(map[storage.CounterName]storage.CounterDateType),
				},
			},
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/value/gauge/metric", nil),
				routerParams: map[string]string{
					"mType": "gauge",
					"mName": "Alloc",
				},
			},
			want: want{
				statusCode: http.StatusOK,
				content:    "10.1",
			},
		},
		{
			name: "test of successfully getting counter metric",
			fields: fields{
				store: &storage.MemStorage{
					Gauges: make(map[storage.GaugeName]storage.GaugeDateType),
					Counters: map[storage.CounterName]storage.CounterDateType{
						"counter1": 1984,
					},
				},
			},
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/value/counter/counter1", nil),
				routerParams: map[string]string{
					"mType": "counter",
					"mName": "counter1",
				},
			},
			want: want{
				statusCode: http.StatusOK,
				content:    "1984",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			for key, value := range tt.args.routerParams {
				rctx.URLParams.Add(key, value)
			}
			request := tt.args.request.WithContext(context.WithValue(tt.args.request.Context(), chi.RouteCtxKey, rctx))

			handler := http.HandlerFunc(NewHandler(tt.fields.store).GetMetricValue)
			handler.ServeHTTP(response, request)
			require.Equal(t, tt.want.statusCode, response.Code)
			if tt.want.statusCode == http.StatusOK {
				assert.Equal(t, tt.want.content, response.Body.String())
			}
		})
	}
}

func TestHandler_GetMetricValueV2(t *testing.T) {
	type fields struct {
		store *storage.MemStorage
	}

	type args struct {
		request *http.Request
	}

	type want struct {
		statusCode int
		content    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "test invalid method",
			fields: fields{store: storage.NewMemStorage()},
			args:   args{request: httptest.NewRequest(http.MethodGet, "/value", nil)},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				content:    "Method not allowed\n",
			},
		},
		{
			name:   "test invalid content",
			fields: fields{store: storage.NewMemStorage()},
			args:   args{request: httptest.NewRequest(http.MethodPost, "/value", strings.NewReader("test"))},
			want: want{
				statusCode: http.StatusBadRequest,
				content:    "invalid character 'e' in literal true (expecting 'r')\n",
			},
		},
		{
			name:   "test - metric not found",
			fields: fields{store: storage.NewMemStorage()},
			args:   args{request: httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(`{"ID":"metric","type":"gauge"}`))},
			want: want{
				statusCode: http.StatusNotFound,
				content:    "Not found.\n",
			},
		},
		{
			name: "test - success get gauge metric",
			fields: fields{
				store: &storage.MemStorage{
					Gauges: map[storage.GaugeName]storage.GaugeDateType{
						"Alloc": 123.123,
					},
					Counters: make(map[storage.CounterName]storage.CounterDateType),
				},
			},
			args: args{request: httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(`{"ID":"Alloc","type":"gauge"}`))},
			want: want{
				statusCode: http.StatusOK,
				content:    `{"id":"Alloc","type":"gauge","value":123.123}`,
			},
		},
		{
			name: "test - success get counter metric",
			fields: fields{
				store: &storage.MemStorage{
					Gauges: make(map[storage.GaugeName]storage.GaugeDateType),
					Counters: map[storage.CounterName]storage.CounterDateType{
						"Counter1": 1984,
					},
				},
			},
			args: args{request: httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(`{"id":"Counter1","type":"counter"}`))},
			want: want{
				statusCode: http.StatusOK,
				content:    `{"id":"Counter1","type":"counter","delta":1984}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response := httptest.NewRecorder()
			NewHandler(tt.fields.store).GetMetricValueV2(response, tt.args.request)
			require.Equal(t, tt.want.statusCode, response.Code)
			assert.Equal(t, tt.want.content, response.Body.String())
		})
	}
}
