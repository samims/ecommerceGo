package handlers

import (
	"context"
	"net/http"
	"strconv"

	protos "currency/protos/currency"
	"product-api/data"
	"product-api/utils"

	"github.com/gorilla/mux"
)

func (p *Products) GetByID(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")

	id := getProductID(r)
	product, err := data.GetProductByID(id)

	if err != nil {
		switch err {
		case data.ErrProductNotFound:
			p.l.Println("[ERROR] fetching product", err)
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			p.l.Println("[ERROR]: fetching product", err)
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// get exchange rate
	rr := protos.RateRequest{
		Base:        protos.Currencies_USD,
		Destination: protos.Currencies_INR,
	}
	resp, err := p.currencyClient.GetRate(context.Background(), &rr)
	if err != nil {
		p.l.Println("[Error] getting new rate", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	p.l.Printf("Currency resp %#v", resp.Rate)

	updatedProduct := *product
	updatedProduct.Price = updatedProduct.Price * resp.Rate
	utils.RespondWithJSON(w, http.StatusOK, updatedProduct)

}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}
