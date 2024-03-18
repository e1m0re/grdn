package apiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name string
		args string
		want *API
	}{
		{
			name: "test API constructor",
			args: "http://localhost:8080",
			want: &API{
				client:  &http.Client{},
				baseURL: "http://localhost:8080",
			},
		},
		// todo Make test with flags???
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(tt.args)
			assert.Equal(t, tt.want, api)
		})
	}
}
