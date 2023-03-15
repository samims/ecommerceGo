package server

import (
	"context"
	"io"
	"time"

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

// SubscribeRates implements the gRPC bidirectional streaming method for the server
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	// handle client messages
	go func() {
		for {
			rr, err := src.Recv() // Recv is a blocking method which returns on client data
			// io.EOF signals that the client has closed the connection
			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}

			// any other error means the transport between the server and client is unavailable
			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				break
			}

			c.log.Info("Handle client request ", "request_base ", rr.GetBase(), " request_dest ", rr.GetDestination())
		}
	}()

	// handle server responses
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// send a message back to the client
			err := src.Send(&protos.RateResponse{Rate: 12.1})
			if err != nil {
				return err
			}
		}
	}
}
