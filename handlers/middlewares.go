package handlers

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"

	"github.com/samims/ecommerceGo/data"
)

type GzipHandler struct {
}

// GZipResponseMiddleWare is a middleware that compresses response using GZip format.
// It takes a http.Handler as an argument, and returns a http.Handler
func (g *GzipHandler) GZipResponseMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrw := NewWrappedResponseWriter(w)
		wrw.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(w, r)
		return
	})
}

// WrappedResponseWriter is a custom http.ResponseWriter that wraps the default http.ResponseWriter
// and uses gzip writer to compress the output.
type WrappedResponseWriter struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
}

// NewWrappedResponseWriter creates a new WrappedResponseWriter, which wraps the default http.ResponseWriter
// and uses gzip writer to compress the output.
func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	// Create a new gzip writer that writes to the underlying http.ResponseWriter.
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{gzipWriter: gw, responseWriter: w}
}

// Header returns the header map that will be sent by WriteHeader.
func (ww *WrappedResponseWriter) Header() http.Header {
	return ww.responseWriter.Header()
}

// Write writes the data to the gzip writer and returns the number of bytes written and any error encountered.
func (ww *WrappedResponseWriter) Write(d []byte) (int, error) {
	return ww.gzipWriter.Write(d)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (ww *WrappedResponseWriter) WriteHeader(statusCode int) {
	ww.responseWriter.WriteHeader(statusCode)
}

type KeyProduct struct{}
type LogVar struct{}

// MiddlewareValidateProduct is a middleware function that validates incoming product data
// by deserializing and validating the product using the Product.Validate() method.
// It returns a http.Handler that can be used in a chain of handlers to handle incoming requests.
// If there is an error during deserialization or validation, the middleware responds with a BadRequest (400) error
// and does not call the next handler in the chain.
//
// Parameters:
// - next: http.Handler - The next handler in the chain to call if validation succeeds.
//
// Returns:
// - http.Handler - An http.Handler that validates incoming product data and calls the next handler in the chain if validation succeeds.
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		// deserialization ...
		err := prod.FromJSON(r.Body)

		if err != nil {
			p.l.Println("[ERROR] deserialization product", err)
			http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}

		// validations ...
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(w, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		// setting the value to context of the request so that we can
		// access it from handlers, of which url we are using the middleware
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

//func (p *Products) LogMiddleWare(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), LogVar{}, log)
//		req := r.WithContext(ctx)
//
//	}
//	}

//func (p *Products) CorsMiddleWare(handler http.Handler) http.Handler {
//	return gorillaHandlers.CORS(
//		gorillaHandlers.AllowedOrigins(p.cfg.AllowedHosts),
//	)(handler)
//}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
