package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestSignChecking(t *testing.T) {
	type args struct {
		method      string
		key         string
		headerName  string
		headerValue string
		body        []byte
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successfully case (GET request)",
			args: args{
				key:         "secret key",
				method:      "GET",
				headerName:  "HashSHA256",
				headerValue: "3fokgzYfs1ICaJVHrp2rNKo03KSMs8uGEfaYL9+AiKA=",
				body:        []byte("request body"),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Request without sum in header",
			args: args{
				key:        "",
				method:     "POST",
				headerName: "",
				body:       make([]byte, 0),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Request without body invalid sum",
			args: args{
				key:         "",
				method:      "POST",
				headerName:  "HashSHA256",
				headerValue: "qwerty",
				body:        make([]byte, 0),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "successfully case (POST request)",
			args: args{
				key:         "secret key",
				method:      "POST",
				headerName:  "HashSHA256",
				headerValue: "LtF60KMiLpS1xdCaUmFOtdvGucz6Y/T+MNI3vtciKHQ=",
				body:        []byte("request body"),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Use(SignChecking(test.args.key))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})

			request := httptest.NewRequest(test.args.method, "/", bytes.NewReader(test.args.body))
			if len(test.args.headerName) > 0 {
				request.Header.Set(test.args.headerName, test.args.headerValue)
			}

			r.ServeHTTP(recorder, request)
			response := recorder.Result()

			require.Equal(t, test.want.statusCode, response.StatusCode)

			response.Body.Close()
		})
	}
}
