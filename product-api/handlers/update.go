package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samims/ecommerceGo/data"
)

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
