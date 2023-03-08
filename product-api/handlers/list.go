package handlers

import (
	"net/http"

	"product-api/utils"
)

// swagger:route GET / productAPIs listProducts
// Returns a list of products
// responses:
//	200: ProductResponseWrapper

// GetProducts retrieves products from the database and returns them in JSON format.
// Args:
//
//	w: an http.ResponseWriter used to write the response to the HTTP client
//	r: a pointer to the http.Request object that contains the request parameters
//
// Returns:
//
//	void
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")

	// Extract the currency from the request URL parameters.
	cr := r.URL.Query().Get("currency")

	// Call the GetProducts method of the product database to retrieve the list of products.
	listProducts, err := p.productDB.GetProducts(cr)

	if err != nil {
		p.l.Error("error getting products")
		utils.RespondWithError(w, http.StatusFailedDependency, "error getting product")
		return
	}

	// Encode the product list as JSON and write it to the response stream
	//err = listProducts.ToJSON(w)
	utils.RespondWithJSON(w, http.StatusOK, listProducts)
	if err != nil {
		http.Error(w, "Unable to marshal", http.StatusInternalServerError)
	}
}
