package compress

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/andybalholm/brotli"
)

// Middleware provides a middleware that compresses the response body using brotli or gzip
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wcloser := brotli.HTTPCompressor(w, r)
		defer wcloser.Close()

		cw := &compressResponseWriter{ResponseWriter: w, compressedWriter: wcloser}
		h.ServeHTTP(cw, r)
	})
}

type compressResponseWriter struct {
	http.ResponseWriter

	// The streaming encoder writer to be used if there is one. Otherwise,
	// this is just the normal writer.
	compressedWriter io.WriteCloser
}

func (cw *compressResponseWriter) Write(p []byte) (int, error) {
	return cw.compressedWriter.Write(p)
}

func (cw *compressResponseWriter) Flush() {
	if f, ok := cw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (cw *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := cw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("http.Hijacker is unavailable on the writer")
}

func (cw *compressResponseWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := cw.ResponseWriter.(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return errors.New("http.Pusher is unavailable on the writer")
}

func (cw *compressResponseWriter) Close() error {
	if err := cw.compressedWriter.Close(); err != nil {
		return err
	}

	if c, ok := cw.ResponseWriter.(io.Closer); ok {
		return c.Close()
	}

	return errors.New("io.Closer is unavailable on the writer")
}

func (cw *compressResponseWriter) Unwrap() http.ResponseWriter {
	return cw.ResponseWriter
}
