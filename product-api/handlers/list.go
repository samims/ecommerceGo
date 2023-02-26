package handlers

import (
	"net/http"

	"product-api/data"
)

// swagger:route GET / productAPIs listProducts
// Returns a list of products
// responses:
//	200: ProductResponseWrapper

// GetProducts ...
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	lp := data.GetProducts()

	// Encode the product list as JSON to make it writable to response stream
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal", http.StatusInternalServerError)
	}
}
