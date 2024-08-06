package middleware

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "successfully case",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		recorder := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Use(Logging())
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

		request := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))

		r.ServeHTTP(recorder, request)
		response := recorder.Result()

		require.Equal(t, test.want.statusCode, response.StatusCode)

		response.Body.Close()
	}
}
