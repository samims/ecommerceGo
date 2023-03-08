// Package handlers Product API
//
// Documentation for Product API
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"product-api/configs"
	"product-api/data"

	"github.com/sirupsen/logrus"
)

//// // A list of products returned in the response
//// swagger:response productsResponse
//type productsResponse struct {
//	// The list of products
//	//
//	// in: body
//	Body []data.Product
//}

//// A list of products returned in the response
//// swagger:response ProductResponseWrapper
//type productsRes struct {
//	// The detail of products
//	//
//	// in: body
//	Body data.A
//}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete
	// in: path
	//required: true
	ID int `json:"id"`
}

// swagger:response noContent
type productsNoContent struct {
}

type Products struct {
	l         *logrus.Logger
	cfg       *configs.Config
	productDB *data.ProductsDB
}

func NewProduct(l *logrus.Logger, pdb *data.ProductsDB) *Products {
	return &Products{
		l:         l,
		productDB: pdb,
	}
}
