package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var ErrProductNotFound = fmt.Errorf("product not found")

// Product defines the structure for API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// require: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

// ProductResponseWrapper is list of product in response
// swagger:response ProductResponseWrapper
type ProductResponseWrapper struct {
	// in: body
	Body []Product
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU, false)

	return validate.Struct(p)

}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)

}

func getNextID() int {
	if len(productList) == 0 {
		return 1
	}
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProducts(id int, p *Product) error {

	idx := findIndexByProductID(id)
	if idx == -1 {
		return ErrProductNotFound
	}
	p.ID = id
	productList[idx] = p
	return nil

}

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func DeleteProduct(id int) error {
	idx := findIndexByProductID(id)
	if idx == -1 {
		return ErrProductNotFound
	}
	productList = append(productList[:idx], productList[idx+1:]...)
	return nil

}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milk cofee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "xyz123",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
}
