package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/samims/ecommerceGO/currency/configs"
	"github.com/samims/ecommerceGO/currency/data"
	"github.com/samims/ecommerceGO/currency/logger"
	"github.com/samims/ecommerceGO/currency/server"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configuration
	os.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "debug")
	cfg := config.NewViperConfig()
	logLevel, err := logrus.ParseLevel(cfg.GetString("log_level"))
	if err != nil {
		fmt.Errorf("error parsing log level: %s", err)
		os.Exit(1)
	}
	log := logger.NewLogger(logLevel)
	// Initialize logger
	if err != nil {
		log.Fatal("Error parsing log level:", err)
	}

	// Initialize server
	s, err := initializeServer(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Error initializing server")
	}

	// Start server
	startServer(s, log)

	// Gracefully shut down server on SIGINT or SIGTERM
	shutdownServer(s, log)
}

func initializeServer(cfg config.Env, log *logrus.Logger) (*server.Server, error) {
	// Initialize rates
	rates, err := data.NewRates(log, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to generate rates: %s", err)
	}

	// Initialize server
	return server.NewServer(cfg, log, rates)
}

func startServer(s *server.Server, log *logrus.Logger) {
	go func() {
		if err := s.Start(); err != nil {
			log.WithError(err).Fatal("Error starting server")
		}
	}()
}

func shutdownServer(s *server.Server, log *logrus.Logger) {
	log.Info("Shutting down server...")

	// Set up a channel to receive the SIGINT or SIGTERM signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Create a context with a timeout of 300 milliseconds
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	ctx := context.Background()

	// Stop the server using the context
	if err := s.Stop(ctx); err != nil {
		log.WithError(err).Fatal("Error stopping server")
	}

	log.Info("Server gracefully stopped")
}
