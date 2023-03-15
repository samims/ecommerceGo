package handlers

import (
	"net/http"
	"strconv"

	"product-api/data"

	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} productAPIs deleteProduct
// Returns blank success
// responses:
//	200: noContent

// DeleteProduct delete a product
func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	p.l.Debugln("Handle delete product", id)

	if err := data.DeleteProduct(id); err != nil {
		switch err {
		case data.ErrProductNotFound:
			p.l.Errorln("Product not found for deletion with id ", id)
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			http.Error(w, "Error deleting product", http.StatusInternalServerError)
		}
		return

	}

}
