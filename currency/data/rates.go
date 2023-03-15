package data

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	config "github.com/samims/ecommerceGO/currency/configs"
	"github.com/samims/ecommerceGO/currency/constants"

	"github.com/sirupsen/logrus"
)

//var rateURI string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRates struct {
	log   *logrus.Logger
	rates map[string]float64
}

func NewRates(l *logrus.Logger, cfg config.Env) (*ExchangeRates, error) {
	exchangeRates := &ExchangeRates{
		log:   l,
		rates: map[string]float64{},
	}
	rateURI := cfg.GetString(constants.EnvRateUri)

	err := exchangeRates.gerRates(rateURI)

	return exchangeRates, err

}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}
	dr, ok := e.rates[dest]

	return dr / br, nil
}

func (e *ExchangeRates) gerRates(uri string) error {
	resp, err := http.DefaultClient.Get(uri)
	if err != nil {
		return fmt.Errorf("Getting error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 got %d", resp.StatusCode)
	}

	defer func(body io.ReadCloser, e *ExchangeRates) {
		if body != nil {
			err := body.Close()
			if err != nil {
				e.log.Error("error closing response body")
			}
		}
	}(resp.Body, e)

	// rest of your code

	cubes := &Cubes{}
	err = xml.NewDecoder(resp.Body).Decode(&cubes)
	if err != nil {
		return err
	}
	for _, c := range cubes.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = r
	}
	e.rates["EUR"] = 1
	return nil

}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
