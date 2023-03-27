package data

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	config "github.com/samims/ecommerceGO/currency/configs"
	"github.com/samims/ecommerceGO/currency/constants"

	"github.com/sirupsen/logrus"
)

//var rateURI string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRates struct {
	log   *logrus.Logger
	mutex *sync.Mutex
	rates map[string]float64
}

func NewRates(l *logrus.Logger, cfg config.Env) (*ExchangeRates, error) {
	exchangeRates := &ExchangeRates{
		log:   l,
		rates: map[string]float64{},
		mutex: &sync.Mutex{},
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

// MonitorRates returns a channel that can be used to monitor currency exchange
func (e *ExchangeRates) MonitorRates(ctx context.Context, interval time.Duration) chan bool {
	// Create a new channel of type struct{} and assign it to ret variable
	ret := make(chan bool)

	// Start a new goroutine that runs the update loop in the background
	go func(interval time.Duration) {
		ticker := time.NewTicker(interval)

		// Enter an infinite loop that waits for messages from the ticker channel
		for {
			select {
			// When a message is received from the ticker channel, simulate currency
			// rate fluctuations and notify any listeners of updates
			case <-ticker.C:
				// Acquire the mutex before accessing e.rates
				e.mutex.Lock()

				// Iterate through each currency rate in the map and modify its value
				// by a random percentage between 10% increase and 10% decrease
				for k, v := range e.rates {
					// Generate a random percentage change between -10% and +10%
					changePercent := rand.Intn(21) - 10 // Returns a random integer between -10 and +10

					// Randomly determine whether the change should be positive or negative
					direction := rand.Intn(2) // Returns either 0 or 1

					// Modify the rate for the current currency by multiplying it with
					// the change percentage and direction
					if direction == 0 {
						// If the direction is negative, subtract the change percentage from the rate
						e.rates[k] = v * (100.0 - float64(changePercent)) / 100
					} else {
						// If the direction is positive, add the change percentage to the rate
						e.rates[k] = v * (100.0 + float64(changePercent)) / 100
					}
				}

				// Release the mutex after we are done accessing e.rates
				e.mutex.Unlock()

				// Notify any listeners of updates by sending an empty struct{} on the ret channel
				// This will block if there is no listener on the other end
				select {
				case ret <- true:
				case <-ctx.Done():
					e.log.Info("Manually shutdown using context1")

					return
				}

			case <-ctx.Done():
				e.log.Info("Manually shutdown using context")
				return
			}
		}
	}(interval)

	// Return the channel to the caller
	return ret
}

// rates at a given time interval
//func (e *ExchangeRates) MonitorRates(ctx context.Context, interval time.Duration) chan bool {
//	// Create a new channel of type struct{} and assign it to ret variable
//	ret := make(chan bool)
//
//	// Start a new goroutine that runs the update loop in the background
//	go func(interval time.Duration) {
//		ticker := time.NewTicker(interval)
//
//		// Enter an infinite loop that waits for messages from the ticker channel
//		for {
//			select {
//			// When a message is received from the ticker channel, simulate currency
//			// rate fluctuations and notify any listeners of updates
//			case <-ticker.C:
//				// Iterate through each currency rate in the map and modify its value
//				// by a random percentage between 10% increase and 10% decrease
//				for k, v := range e.rates {
//					// Generate a random percentage change between -10% and +10%
//					changePercent := rand.Intn(21) - 10 // Returns a random integer between -10 and +10
//
//					// Randomly determine whether the change should be positive or negative
//					direction := rand.Intn(2) // Returns either 0 or 1
//
//					// Modify the rate for the current currency by multiplying it with
//					// the change percentage and direction
//					if direction == 0 {
//						// If the direction is negative, subtract the change percentage from the rate
//						e.rates[k] = v * (100.0 - float64(changePercent)) / 100
//					} else {
//						// If the direction is positive, add the change percentage to the rate
//						e.rates[k] = v * (100.0 + float64(changePercent)) / 100
//					}
//				}
//
//				// Notify any listeners of updates by sending an empty struct{} on the ret channel
//				// This will block if there is no listener on the other end
//				select {
//				case ret <- true:
//				case <-ctx.Done():
//					e.log.Info("Manually shutdown using context1")
//
//					return
//				}
//
//			case <-ctx.Done():
//				e.log.Info("Manually shutdown using context")
//				return
//			}
//		}
//	}(interval)
//
//	// Return the channel to the caller
//	return ret
//}

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
