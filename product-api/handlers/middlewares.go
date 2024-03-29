package handlers

import (
	"context"
	"fmt"
	"net/http"

	"product-api/data"
)

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
			p.l.Errorln("unable to deserialize product", err)
			http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}

		// validations ...
		err = prod.Validate()
		if err != nil {
			p.l.Errorln("unable validating product", err)
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

// CORS  is custom cors middleware
func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: write your allowed origin here, you can little customize
		// & pass the origin through dependency injection
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
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
