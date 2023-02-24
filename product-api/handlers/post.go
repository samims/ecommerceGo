package handlers

import (
	"net/http"

	"product-api/data"
)

// swagger:route POST /products productAPIs createProduct
//

// Create adds a new product to the product list.
// The product information is extracted from the context of the request.
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {
	// fetch the data from context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %v\n", prod)

	data.AddProduct(&prod)
	w.WriteHeader(http.StatusCreated)

}
