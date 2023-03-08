package server

import (
	"context"

	"github.com/samims/ecommerceGO/currency/data"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"

	"github.com/sirupsen/logrus"
)

// Currency represents a server that provides currency conversion rates.
type Currency struct {
	log *logrus.Logger
	protos.UnimplementedCurrencyServer
	rates *data.ExchangeRates
}

func NewCurrency(l *logrus.Logger, r *data.ExchangeRates) *Currency {
	return &Currency{
		log:   l,
		rates: r,
	}
}

// GetRate retrieves the exchange rate for the given currencies.
//
// ctx: The context for the request.
// rr: The RateRequest, containing the base and destination currencies.
//
// Returns a RateResponse containing the exchange rate.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate ", " base ", rr.GetBase(), " destination ", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.Base.String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: rate}, err
}
