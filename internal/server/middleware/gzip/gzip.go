package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, err
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {

			contentEncoding := request.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")

			if sendsGzip {
				cr, err := newCompressReader(request.Body)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}

				request.Body = cr
				defer cr.Close()
			}

			next.ServeHTTP(writer, request)
		},
	)
}
