package handlers

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/samims/ecommerceGO/currency/data"
	pb "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	contextClientIDKey contextKey = iota
)

// CurrencyService represents a handlers that provides currency conversion rates.
type CurrencyService struct {
	log           *logrus.Logger
	ctx           context.Context
	rates         *data.ExchangeRates
	subscriptions map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest
	pb.UnimplementedCurrencyServer
}

// NewCurrency creates a new instance of the CurrencyService with the given context, logger, and exchange rates.
// It initializes the subscriptions and clients maps and starts a goroutine to handle rate updates.
// It returns a pointer to the *CurrencyService instance.
func NewCurrency(ctx context.Context, l *logrus.Logger, r *data.ExchangeRates) *CurrencyService {
	c := &CurrencyService{
		ctx:           ctx,
		log:           l,
		rates:         r,
		subscriptions: make(map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest),
	}

	go c.handleUpdates()

	return c
}

// SubscribeRates is a gRPC streaming method that allows clients to subscribe to currency rate updates.
// It listens for incoming client requests and sends updates periodically based on a ticker.
// Each subscribed client is added to a subscription list, which is used to send updates to all subscribed clients.
func (c *CurrencyService) SubscribeRates(stream pb.Currency_SubscribeRatesServer) error {
	// generate a unique client ID
	clientID := getClientID(c.ctx)

	// register the client
	//c.clients[clientID] = stream

	// notify that a new client has connected
	c.log.Infof("client %s connected", clientID)

	// start an infinite loop to send rate updates
	for {
		// read the request from the client
		req, err := stream.Recv()
		if err == io.EOF {
			// client has disconnected
			c.log.Infof("client %s disconnected", clientID)
			break
		}
		if err != nil {
			c.log.Errorf("error receiving stream from client %s: %v", clientID, err)
			return err
		}

		// log the request
		c.log.Infof("Handling client request.%s", req.String())

		rrs, ok := c.subscriptions[stream]
		if !ok {
			rrs = []*pb.RateRequest{}
		}

		rrs = append(rrs, req)
		c.subscriptions[stream] = rrs
	}

	return nil
}

// handleUpdates sends updated currency exchange rate to subscribed clients
// every 5 seconds.
func (c *CurrencyService) handleUpdates() {
	// Monitor exchange rate updates using a ticker with a 5-second interval
	rateUpdates := c.rates.MonitorRates(c.ctx, 5*time.Second)

	// Continuously loop over the ticker channel to receive rate updates
	for range rateUpdates {
		c.log.Info("got updated rates")
		// Update all subscribed clients with the new exchange rate
		c.updateSubscriptions()
	}
}

// updateSubscriptions sends the updated currency exchange rate to each
// subscribed client.
func (c *CurrencyService) updateSubscriptions() {
	// Loop over all subscribed clients and their requested currency pairs
	for k, v := range c.subscriptions {
		for _, rr := range v {
			// Get the updated exchange rate for the client's currency pair
			rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
			if err != nil {
				// Log an error message if the exchange rate could not be retrieved
				c.log.Error(
					"unable to get updated rate",
					"base", rr.GetBase(),
					"destination", rr.GetDestination(),
				)
			}
			err = k.Send(&pb.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate})
		}
	}
}

// GetRate retrieves the exchange rate for the given base and destination currencies.
// It returns a RateResponse containing the exchange rate and the base and destination currencies.
func (c *CurrencyService) GetRate(_ context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate ", " base ", rr.GetBase(), " destination ", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.Base.String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	rateResp := &pb.RateResponse{
		Rate:        rate,
		Base:        rr.GetBase(),
		Destination: rr.GetDestination()}
	return rateResp, nil

}

func getClientID(ctx context.Context) string {
	id, ok := ctx.Value(contextClientIDKey).(string)
	if !ok {
		// Generate a new UUID as client ID
		uuid := uuid.New().String()
		// Add the UUID to the context
		ctx = context.WithValue(ctx, contextClientIDKey, uuid)
		return uuid
	}
	return id
}
