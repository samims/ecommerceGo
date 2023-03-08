package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	protos "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/sirupsen/logrus"

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
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      *logrus.Logger
}

func NewProductsDB(c protos.CurrencyClient, l *logrus.Logger) *ProductsDB {
	return &ProductsDB{
		currency: c,
		log:      l,
	}
}

// ProductResponseWrapper is list of product in response
// swagger:response ProductResponseWrapper
type ProductResponseWrapper struct {
	// in: body
	Body []Product
}

func (p *Product) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("sku", validateSKU, false)
	if err != nil {
		return err
	}

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

// GetProducts retrieves a list of products in a given currency from the ProductsDB.
// If currency is empty, returns the original product list. Otherwise, it gets the
// exchange rate using getRate and applies it to each product's price to create a new
// list. May modify the ProductsDB state if getRate is called.
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	// If the currency parameter is empty, return the original productList.
	if currency == "" {
		return productList, nil
	}

	// Otherwise, retrieve the exchange rate for the given currency.
	rate, err := p.getRate(currency)

	if err != nil {
		p.log.Error("unable to get rate currency", currency, "error", err)
		return nil, err
	}

	// Create a new Products slice to avoid modifying the original
	// productList, used in different places
	pr := Products{}

	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}
	return pr, nil
}

// GetProductByID retrieves a product with a given ID from the ProductsDB.
// If the product is not found, returns an error.
// If currency is empty, returns the product.
// Otherwise, it gets the exchange rate using getRate and applies it to the product's price.
// May modify the ProductsDB state if getRate is called.
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	idx := findIndexByProductID(id)
	if idx == -1 {
		return new(Product), ErrProductNotFound
	}
	if currency == "" {
		return productList[idx], nil
	}
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("unable to get rate", currency, currency, "error", err)
		return nil, err
	}

	// This is done to avoid modifying the original product list.
	npObj := *productList[idx]
	npObj.Price = npObj.Price * rate
	return &npObj, nil
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

func (p *ProductsDB) UpdateProducts(id int, pObj *Product) error {

	idx := findIndexByProductID(id)
	if idx == -1 {
		return ErrProductNotFound
	}
	pObj.ID = id
	productList[idx] = pObj
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

func (p *ProductsDB) getRate(destination string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rr)
	return resp.Rate, err
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
		Description: "Frothy milk coffee",
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
