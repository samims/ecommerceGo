package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samims/ecommerceGo/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l: l}
}

// GetProducts retrieves a list of all products from the database and writes it as JSON to the response.
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	lp := data.GetProducts()

	// Encode the product list as JSON to make it writable to response stream
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal", http.StatusInternalServerError)
	}
}

// AddProduct adds a new product to the product list.
// The product information is extracted from the context of the request.
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)

}

// UpdateProducts update the product with the given ID using data from the request body.
// The product ID is extracted from the request URL.
func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {

	// Extract the product ID from the URL.
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "unable to convert id ", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle put request", id)

	// Extract the product data from the request context.
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	// Update the product with the specified ID.
	if err := data.UpdateProducts(id, &prod); err != nil {
		switch err {
		case data.ErrProductNotFound:
			p.l.Println("[ERROR] product not found with provided id ", id)
			http.Error(w, "Product not found", http.StatusNotFound)
		default:
			http.Error(w, "Error updating product", http.StatusInternalServerError)
		}
		return
	}
}

type KeyProduct struct{}

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

		next.ServeHTTP(w, req)
	})
}
