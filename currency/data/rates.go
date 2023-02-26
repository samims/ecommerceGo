package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var rateURI string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRates struct {
	log   *logrus.Logger
	rates map[string]float64
}

func NewRates(l *logrus.Logger) (*ExchangeRates, error) {
	exchangeRates := &ExchangeRates{
		log:   l,
		rates: map[string]float64{},
	}

	err := exchangeRates.gerRates(rateURI)

	return exchangeRates, err

}

func (e *ExchangeRates) gerRates(uri string) error {
	resp, err := http.DefaultClient.Get(uri)
	if err != nil {

	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

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
	return nil

}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
