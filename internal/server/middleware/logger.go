package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size

	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					slog.Error("Internal error",
						slog.String("error", fmt.Sprintf("%v", err)),
						slog.String("stack", string(debug.Stack())),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(writer)
			next.ServeHTTP(wrapped, request)

			duration := time.Since(start)
			slog.Info("Incoming request",
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", wrapped.status),
				slog.Duration("duration", duration),
				slog.Int("size", wrapped.size),
			)
		},
	)
}
