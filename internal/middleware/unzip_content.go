package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// UnzipContent extracts requests body.
func UnzipContent() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				contentEncoding := r.Header.Get("Content-Encoding")

				if strings.Contains(contentEncoding, "gzip") {
					gr, err := gzip.NewReader(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					r.Body = gr
				}

				next.ServeHTTP(w, r)
			})
	}
}
