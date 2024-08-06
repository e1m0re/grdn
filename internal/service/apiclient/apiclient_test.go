package apiclient

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAPIClient_DoRequest(t *testing.T) {
	type fields struct {
		testServer *httptest.Server
	}
	type args struct {
		request func(u *url.URL, s *httptest.Server) *http.Request
	}
	type want struct {
		response *http.Response
		err      error
	}
	tests := []struct {
		args   args
		fields fields
		want   want
		name   string
	}{
		{
			name: "Response 500",
			fields: fields{
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})),
			},
			args: args{
				request: func(u *url.URL, s *httptest.Server) *http.Request {
					return &http.Request{
						Method: u.Scheme,
						URL:    u,
						Header: make(http.Header),
					}
				},
			},
			want: want{
				response: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				},
				err: nil,
			},
		},
		{
			name: "Response 404",
			fields: fields{
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				})),
			},
			args: args{
				request: func(u *url.URL, s *httptest.Server) *http.Request {
					return &http.Request{
						Method: u.Scheme,
						URL:    u,
						Header: make(http.Header),
					}
				},
			},
			want: want{
				response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       http.NoBody,
				},
				err: nil,
			},
		},
		{
			name: "Response 200 with empty body",
			fields: fields{
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})),
			},
			args: args{
				request: func(u *url.URL, s *httptest.Server) *http.Request {
					return &http.Request{
						Method: u.Scheme,
						URL:    u,
						Header: make(http.Header),
					}
				},
			},
			want: want{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       http.NoBody,
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() { test.fields.testServer.Close() }()
			u, _ := url.Parse(test.fields.testServer.URL)
			apiClient := NewAPIClient(u.Scheme, u.Host, nil)
			got, err := apiClient.DoRequest(test.args.request(u, test.fields.testServer))
			if got != nil {
				defer got.Body.Close()
			}
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.response.StatusCode, got.StatusCode)
		})
	}
}

func TestAPIClient_SendMetricsData(t *testing.T) {
	type fields struct {
		testServer *httptest.Server
		key        []byte
	}
	type args struct {
		data []byte
	}
	type want struct {
		err error
	}
	tests := []struct {
		args   args
		want   want
		name   string
		fields fields
	}{
		{
			name: "server response with code 500",
			fields: fields{
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})),
			},
			args: args{data: make([]byte, 0)},
			want: want{err: nil},
		},
		{
			name: "Successfully case with encrypt",
			fields: fields{
				key: []byte("key"),
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})),
			},
			args: args{data: make([]byte, 0)},
			want: want{err: nil},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() { test.fields.testServer.Close() }()
			u, _ := url.Parse(test.fields.testServer.URL)
			apiClient := NewAPIClient("http", u.Host, test.fields.key)
			err := apiClient.SendMetricsData(&test.args.data)
			assert.Equal(t, test.want.err, err)
		})
	}
}
