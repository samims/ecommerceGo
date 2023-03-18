package handlers

import (
	"context"
	"fmt"
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
	log *logrus.Logger
	pb.UnimplementedCurrencyServer
	ctx           context.Context
	rates         *data.ExchangeRates
	subscriptions map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest
	clients       map[string]pb.Currency_SubscribeRatesServer
}

func NewCurrency(ctx context.Context, l *logrus.Logger, r *data.ExchangeRates) *CurrencyService {
	c := &CurrencyService{
		ctx:           ctx,
		log:           l,
		rates:         r,
		subscriptions: make(map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest),
		clients:       make(map[string]pb.Currency_SubscribeRatesServer),
	}

	go c.handleUpdates()

	return c
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
			} else {
				// Send the updated exchange rate to the subscribed client
				c.sendRateResponse(k, rr, rate)
			}
		}
	}
}

// sendRateResponse sends the updated currency exchange rate to a subscribed client.
func (c *CurrencyService) sendRateResponse(stream pb.Currency_SubscribeRatesServer, rr *pb.RateRequest, rate float64) {
	// Send the updated exchange rate to the client's stream
	resp := &pb.RateResponse{
		Base:        rr.GetBase(),
		Destination: rr.GetDestination(),
		Rate:        rate,
	}
	err := stream.Send(resp)
	if err != nil {
		// Log an error message if the exchange rate update could not be sent to the client
		c.log.Error("unable to send rate update to client", "error", err)
	}
}

//func (c *CurrencyService) handleUpdates() {
//	rateUpdates := c.rates.MonitorRates(5 * time.Second)
//	for range rateUpdates {
//		c.log.Info("got updated rates")
//		for k, v := range c.subscriptions {
//			for _, rr := range v {
//				rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
//				if err != nil {
//					c.log.Error(
//						"unable to get updated rate",
//						"base", rr.GetBase(),
//						"destination", rr.GetDestination(),
//					)
//				}
//				k.Send(&pb.RateResponse{
//					Base:        rr.GetBase(),
//					Destination: rr.GetDestination(),
//					Rate:        rate,
//				})
//			}
//		}
//	}
//
//}

// GetRate retrieves the exchange rate for the given currencies.
//
// ctx: The context for the request.
// rr: The RateRequest, containing the base and destination currencies.
//
// Returns a RateResponse containing the exchange rate.
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

func (c *CurrencyService) SubscribeRates(stream pb.Currency_SubscribeRatesServer) error {
	// generate a unique client ID
	clientID := uuid.NewString()

	// register the client
	c.clients[clientID] = stream

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

		// send rate updates to the client
		err = c.sendRateUpdates(req, clientID, stream)
		//c.updateSubscriptions()

		if err != nil {
			return err
		}
	}

	// unregister the client
	delete(c.clients, clientID)

	return nil
}

func (c *CurrencyService) sendRateUpdates(req *pb.RateRequest, clientID string, stream pb.Currency_SubscribeRatesServer) error {
	for {
		// get the rate for the requested currency pair
		rate, err := c.rates.GetRate(req.GetBase().String(), req.GetDestination().String())
		if err != nil {
			return err
		}

		// create a rate update message
		update := &pb.RateResponse{
			Base:        req.GetBase(),
			Destination: req.GetDestination(),
			Rate:        rate,
		}

		// send the rate update to the client
		err = stream.Send(update)
		if err != nil {
			c.log.Errorf("error sending rate update to client %s: %v", clientID, err)
			return err
		}

		// wait for a second before sending the next update
		time.Sleep(5 * time.Second)
	}
}

func getClientID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextClientIDKey).(string)
	if !ok {
		// Generate a new UUID as client ID
		uuid := uuid.New().String()
		// Add the UUID to the context
		ctx = context.WithValue(ctx, contextClientIDKey, uuid)
		return uuid, true
	}
	return id, true
}

func (c *CurrencyService) sendRateUpdate(req *pb.RateRequest, clientID string) error {
	rate, err := c.rates.GetRate(req.Base.String(), req.Destination.String())
	if err != nil {
		c.log.Errorf("error getting exchange rate: %v", err)
		return err
	}

	response := &pb.RateResponse{
		Base:        req.Base,
		Destination: req.Destination,
		Rate:        rate,
	}

	if stream, ok := c.clients[clientID]; ok {
		err = stream.Send(response)
		if err != nil {
			c.log.Errorf("unable to send rate update to client %s: %v", clientID, err)
			return err
		}
	} else {
		c.log.Errorf("unable to send rate update to client %s: stream not found", clientID)
		return fmt.Errorf("stream not found for client %s", clientID)
	}

	c.log.Infof("sent rate update to client %s: %v", clientID, response)
	return nil
}

//func (c *CurrencyService) SubscribeRates(stream pb.Currency_SubscribeRatesServer) error {
//	// Loop until the client disconnects or an error occurs
//	for {
//		rr, err := stream.Recv()
//		if err == io.EOF {
//			c.log.Info("client has closed connection")
//			return nil
//		}
//		if err != nil {
//			c.log.Error("unable to read from client ", "error", err)
//			return err
//		}
//
//		c.log.Info("Handle client req ", "req_base", rr.GetBase(), "req_dest", rr.GetDestination())
//
//		// Append the new subscription to the existing list of subscriptions
//		c.subscriptions[stream] = append(c.subscriptions[stream], rr)
//
//		// Send the initial currency rate to the client
//		rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
//		if err != nil {
//			c.log.Error("unable to get rate ", "error", err)
//			return err
//		}
//		if err := stream.Send(&pb.RateResponse{Rate: rate}); err != nil {
//			c.log.Error("unable to send rate to client ", "error", err)
//			return err
//		}
//
//		// Create a ticker to send updates at a regular interval
//		ticker := time.NewTicker(time.Second * 10)
//		defer ticker.Stop()
//
//		// Continuously send updates to the client for the subscribed currencies
//		for range ticker.C {
//			rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
//			if err != nil {
//				c.log.Error("unable to get rate ", "error", err)
//				return err
//			}
//			if err := stream.Send(&pb.RateResponse{Rate: rate}); err != nil {
//				c.log.Error("unable to send rate to client ", "error", err)
//				return err
//			}
//		}
//	}
//}

// SubscribeRates streams currency exchange rate updates to the client over the
// given gRPC stream. The function blocks until the client disconnects or an error
// occurs. During the stream, any rate requests received from the client are added
// to the list of active subscriptions.
//
// Parameters:
//
//	stream: A gRPC stream over which to stream rate updates to the client.
//
// Returns:
//
//	An error if there was a problem reading from the client stream.

//	func (c *CurrencyService) xxSubscribeRates(stream pb.Currency_SubscribeRatesServer) error {
//		clientID, ok := getClientID(c.ctx)
//		if !ok {
//			c.log.Error("client ID not found in context")
//			return errors.New("client ID not found in context")
//		}
//		c.clients[clientID] = stream
//		defer delete(c.clients, clientID)
//
//		c.log.Infof("client %s connected", clientID)
//
//		for {
//			req, err := stream.Recv()
//			if err == io.EOF {
//				c.log.Infof("client %s disconnected", clientID)
//				return nil
//			}
//			if err != nil {
//				c.log.Errorf("error receiving stream from client %s: %v", clientID, err)
//				return err
//			}
//			c.log.Infof("Handling client request.%s", req.String())
//
//			err = c.sendRateUpdate(req, clientID)
//			if err != nil {
//				return err
//			}
//		}
//	}

func (c *CurrencyService) SubscribeRatess(stream pb.Currency_SubscribeRatesServer) error {
	// Loop until the client disconnects or an error occurs
	for {
		rr, err := stream.Recv()
		c.log.Info("xx ", rr.GetBase().String())
		if err == io.EOF {
			c.log.Info("client has closed connection")
			return nil
		}
		if err != nil {
			c.log.Error("unable to read from client ", "error", err)
			return err
		}

		c.log.Info("Handle client req ", "req_base ", rr.GetBase(), " req_dest ", rr.GetDestination())

		// Append the new subscription to the existing list of subscriptions
		c.subscriptions[stream] = append(c.subscriptions[stream], rr)
	}
}
