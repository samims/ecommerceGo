package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

// GzipMiddleware is a middleware that compresses response using GZip format.
// It takes a http.Handler as an argument, and returns a http.Handler
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			// create a gzipped response
			wrw := NewWrappedResponseWriter(writer)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, request)
			defer wrw.Flush()
			return
		}
		next.ServeHTTP(writer, request)
	})
}

// WrappedResponseWriter is a custom http.ResponseWriter that wraps the default http.ResponseWriter
// and uses gzip writer to compress the output.
type WrappedResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

// NewWrappedResponseWriter creates a new WrappedResponseWriter, which wraps the default http.ResponseWriter
// and uses gzip writer to compress the output.
func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{
		w:  w,
		gw: gw,
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.w.Header()
}

// Write writes the data to the gzip writer and returns the number of bytes written and any error encountered.
func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.w.WriteHeader(statusCode)
}

func (wr *WrappedResponseWriter) Flush() {
	err := wr.gw.Flush()
	if err != nil {
		return
	}
	err = wr.gw.Close()
	if err != nil {
		return
	}
}
