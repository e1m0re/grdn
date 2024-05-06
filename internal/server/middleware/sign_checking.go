package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

func SignChecking(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			ctrlSum := r.Header.Get("HashSHA256")
			if ctrlSum == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			h := hmac.New(sha256.New, []byte(key))
			h.Write(body)
			sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
			if sum != ctrlSum {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
