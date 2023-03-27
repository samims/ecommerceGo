package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/samims/ecommerceGO/currency/configs"
	"github.com/samims/ecommerceGO/currency/constants"
	"github.com/samims/ecommerceGO/currency/data"
	"github.com/samims/ecommerceGO/currency/handlers"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg   config.Env
	log   *logrus.Logger
	rates *data.ExchangeRates
	gs    *grpc.Server
}

func NewServer(cfg config.Env, log *logrus.Logger, rates *data.ExchangeRates) (*Server, error) {

	gs := grpc.NewServer()

	cs := handlers.NewCurrency(context.Background(), log, rates)
	protos.RegisterCurrencyServer(gs, cs)
	reflection.Register(gs)

	return &Server{
		cfg:   cfg,
		log:   log,
		rates: rates,
		gs:    gs,
	}, nil
}

func (s *Server) Start() error {
	portStr := s.cfg.GetString(constants.EnvPort)
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", portStr))
	if err != nil {
		return err
	}
	s.log.Info("Serving on port ", portStr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		if err := s.gs.Serve(l); err != nil {
			s.log.Fatal("unable to serve grpc ", err)
		}
	}()

	<-sigChan

	return nil
}

//func (s *Server) Stop(ctx context.Context) error {
//	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//
//	s.gs.GracefulStop()
//
//	select {
//	case <-ctx.Done():
//		s.log.Warn("Shutdown timed out")
//	case <-time.After(500 * time.Millisecond):
//		s.log.Warn("Shutdown timed out")
//	}
//
//	return nil
//}

//func (s *Server) Stop(ctx context.Context) error {
//	//ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
//	//defer cancel()
//	s.log.Info("stopping the server")
//
//	s.gs.GracefulStop()
//	//s.rates.Stop()
//
//	select {
//	case <-ctx.Done():
//		s.log.Warn("Shutdown timed out")
//	case <-time.After(500 * time.Millisecond):
//		s.log.Warn("Shutdown timed out")
//	}
//
//	s.log.Info("Server stopped")
//	return nil
//}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping the server")

	// Stop the gRPC server
	s.gs.Stop()

	// Wait for up to 5 seconds for all connections to be closed
	timeout := 5 * time.Second
	done := make(chan struct{})
	go func() {
		s.gs.GracefulStop()
		close(done)
	}()

	// Wait for the server to stop or the context to be cancelled
	select {
	case <-ctx.Done():
		s.log.Warn("Shutdown timed out")
	case <-done:
		s.log.Info("Server stopped")
	case <-time.After(timeout):
		s.log.Warn("Shutdown timed out")
	}

	return nil
}
