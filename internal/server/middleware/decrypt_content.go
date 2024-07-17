package middleware

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/e1m0re/grdn/internal/encryption"
)

type decryptReader struct {
	r         io.ReadCloser
	decryptor encryption.Decryptor
}

func newDecryptReader(r io.ReadCloser, decryptor encryption.Decryptor) *decryptReader {
	return &decryptReader{
		r:         r,
		decryptor: decryptor,
	}
}

func (c *decryptReader) Read(p []byte) (n int, err error) {
	p, err = c.decryptor.Decrypt(p)

	return len(p), err
}

func (c *decryptReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.Close()
}

// DecryptContent decrypts requests body.
func DecryptContent(privateKeyFile string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					next.ServeHTTP(w, r)
					return
				}

				decryptor, err := encryption.NewDecryptor(privateKeyFile)
				if err != nil {
					slog.Error("internal server error", slog.String("error", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				dr := newDecryptReader(r.Body, decryptor)
				r.Body = dr
				defer dr.Close()

				next.ServeHTTP(w, r)
			})
	}
}
