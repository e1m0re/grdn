package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestSignResponse(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		headerName    string
		headerContent string
		statusCode    int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successfully case",
			args: args{
				key: "secret key",
			},
			want: want{
				statusCode:    200,
				headerName:    "HashSHA256",
				headerContent: "I5/FHTlJaYQFYx9mBuu5XcBOf8aVGxxUGK9GHnV4dZo=",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Use(SignResponse(test.args.key))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

			request := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))

			r.ServeHTTP(recorder, request)
			response := recorder.Result()

			require.Equal(t, test.want.statusCode, response.StatusCode)
			require.Equal(t, test.want.headerContent, response.Header.Get(test.want.headerName))

			response.Body.Close()
		})
	}
}
