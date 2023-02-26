package server

import (
	"context"

	protos "currency/protos/currency"

	"github.com/sirupsen/logrus"
)

type Currency struct {
	log *logrus.Logger
	protos.UnimplementedCurrencyServer
}

//func (c *Currency) mustEmbedUnimplementedCurrencyServer() {
//	//TODO implement me
//	panic("implement me")
//}

func NewCurrency(l *logrus.Logger) *Currency {
	return &Currency{log: l}

}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate ", " base ", rr.GetBase(), " destination ", rr.GetDestination())
	return &protos.RateResponse{
		Rate: 0.5,
	}, nil
}
