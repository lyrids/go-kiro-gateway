package middleware

import (
	"io"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// countingReader wraps a request body and counts the bytes read from it.
type countingReader struct {
	io.ReadCloser
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.ReadCloser.Read(p)
	c.n += int64(n)
	return n, err
}

// RequestLogger returns a middleware that logs one line per request with the
// method, path, status, duration, and request/response body sizes in bytes.
//
// request_bytes counts what the handlers actually read, so it is accurate for
// chunked bodies (where Content-Length is -1). response_bytes counts what was
// written, including streamed SSE output. The chi wrapper keeps http.Flusher
// working for streaming handlers.
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			body := &countingReader{ReadCloser: r.Body}
			if r.Body != nil {
				r.Body = body
			}
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK // handler wrote nothing
			}
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", status).
				Int64("request_bytes", body.n).
				Int64("content_length", r.ContentLength).
				Int("response_bytes", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Msgf("HTTP %d - %s %s", status, r.Method, r.URL.Path)
		})
	}
}
