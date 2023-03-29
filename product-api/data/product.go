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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l *logrus.Logger) *ProductsDB {
	pdb := &ProductsDB{
		currency: c,
		log:      l,
		rates:    make(map[string]float64),
	}

	go pdb.handleUpdates()

	return pdb
}

func (p *ProductsDB) handleUpdates() {
	subscribedClient, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("unable to subscribe for rates ", " error ", err)
	}
	p.client = subscribedClient

	for {
		rr, err := subscribedClient.Recv()
		if err == io.EOF {
			p.log.Error("eof receiving message ", " error ", err)
			break
		}
		if err != nil {
			p.log.Errorf("error receiving message: %v", err)
			continue
		}
		if grpcError := rr.GetError(); grpcError != nil {
			p.log.Errorf("error subscribing for rate  error: %s", err)
			continue
		}
		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("received updated rate ", " dest ", resp.Destination.String(), " ", resp.Rate)
			p.rates[rr.GetRateResponse().Destination.String()] = resp.Rate
		}
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
// Parameters:
//
//	id (int): The ID of the product to retrieve.
//	currency (string): The currency to use for the product price (optional).
//
// Returns:
//
//	*Product: The product object.
//	error: Returns an error if the product is not found or if an error occurred.
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
// May modify the ProductsDB state if getRate is called.
//
// Parameters:
// id (int): The ID of the product to retrieve.
// currency (string): The currency in which to retrieve the product's price.
//
// Returns:
// (*Product): A pointer to the retrieved product object.
// error: Returns an error if the product is not found or if there's an issue with the currency rate conversion.
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

// UpdateProducts updates a product in the database by ID.
// Parameters:
//
//	id (int): The ID of the product to update.
//	pObj (*Product): The updated product object.
//
// Returns:
//
//	error: Returns an error if the product is not found.
func (p *ProductsDB) UpdateProducts(id int, pObj *Product) error {

	idx := findIndexByProductID(id)
	if idx == -1 {
		return ErrProductNotFound
	}
	pObj.ID = id
	productList[idx] = pObj
	return nil

}

// findIndexByProductID searches the productList slice for a product with the given ID and returns the index of the first
// matching product. If no product with the given ID is found, it returns -1.
// Args:
//   id (int): The ID of the product to search for.
// Returns:
//   int: The index of the first product in the productList slice with the given ID, or -1 if no such product is found.

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// getRate returns the current exchange rate for the given destination currency. If the rate is not cached,
// the function will call the Currency service to get the current rate. The rate will be cached for future
// calls to improve performance. The function also sends a subscription request to the Currency service to receive
// updates on the exchange rate for the given destination currency.
//
// Parameters:
//   - destination (string): The 3-letter currency code of the destination currency for which the exchange rate is requested.
//
// Returns:
//   - float64: the exchange rate for the given destination currency.
//   - error: an error indicating whether the request to the currency service was successful
//     or not.
func (p *ProductsDB) getRate(destination string) (float64, error) {
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}
	return p.fetchRate(destination)
}

func (p *ProductsDB) fetchRate(destination string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			md := s.Details()[0].(*protos.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf(
					"unable to get rate from currency server, base - %s & destination - %s is same ",
					md.Base.String(),
					md.Destination.String(),
				)
			}
			return -1, fmt.Errorf(
				"%s, base %s & destination %s",
				s.Err().Error(),
				md.Base.String(),
				md.Destination.String(),
			)
		}
		return -1, err
	}

	p.rates[destination] = resp.Rate
	p.client.Send(rr)

	return resp.Rate, nil
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
