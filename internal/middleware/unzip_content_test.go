package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnzipContent(t *testing.T) {
	type args struct {
		body           []byte
		compressedBody bool
	}
	type want struct {
		body       []byte
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successfully case (without compress)",
			args: args{
				body:           []byte("request body"),
				compressedBody: false,
			},
			want: want{
				body:       []byte("request body"),
				statusCode: 200,
			},
		},
		{
			name: "successfully case (with compress)",
			args: args{
				body:           []byte("request body"),
				compressedBody: true,
			},
			want: want{
				body:       []byte("request body"),
				statusCode: 200,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Use(UnzipContent())
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				if _, err := w.Write([]byte(b)); err != nil {
					panic(err)
				}
			})

			bodyReader := bytes.NewReader(test.args.body)
			if test.args.compressedBody {
				var buf bytes.Buffer
				gzipWriter := gzip.NewWriter(&buf)
				if _, err := gzipWriter.Write(test.args.body); err != nil {
					panic(err)
				}
				gzipWriter.Close()
				bodyReader = bytes.NewReader(buf.Bytes())
			}
			request := httptest.NewRequest("POST", "/", bodyReader)
			request.Close = true
			if test.args.compressedBody {
				request.Header.Set("Content-Encoding", "gzip")
			}

			r.ServeHTTP(recorder, request)
			response := recorder.Result()

			require.Equal(t, test.want.statusCode, response.StatusCode)
			assert.Equal(t, test.want.body, recorder.Body.Bytes())

			response.Body.Close()
		})
	}
}
