package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

// SignResponse signs server responses.
func SignResponse(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := hmac.New(sha256.New, []byte(key))

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Write([]byte(r.URL.Path))
			w.Header().Set("HashSHA256", base64.StdEncoding.EncodeToString(h.Sum(nil)))
			next.ServeHTTP(w, r)
		})
	}
}
