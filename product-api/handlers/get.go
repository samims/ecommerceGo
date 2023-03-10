package handlers

import (
	"net/http"
	"strconv"

	"product-api/data"
	"product-api/utils"

	"github.com/gorilla/mux"
)

func (p *Products) GetByID(w http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get record")
	r.Header.Add("Content-Type", "application/json")

	cur := r.URL.Query().Get("currency")
	id := getProductID(r)
	product, err := p.productDB.GetProductByID(id, cur)

	if err != nil {
		switch err {
		case data.ErrProductNotFound:
			p.l.Errorf("unable fetching product %s", err.Error())
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			p.l.Errorf("unable fetching product %s", err.Error())
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, product)

}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}
