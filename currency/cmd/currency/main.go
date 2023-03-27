package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

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
	startServer(s, log)
	// Start the server
	// This function will start the server in a new goroutine
	// and will not block the main goroutine
	// It will log any errors encountered during the startup process	startServer(s, log)
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
	log.Info("Starting the server..")
	go func() {
		if err := s.Start(); err != nil {
			log.WithError(err).Fatal("Error starting server")
		}
	}()
}

func shutdownServer(s *server.Server, log *logrus.Logger) {
	// Set up a channel to receive the OS interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Wait for the interrupt signal
	<-sigChan
	log.Info("Received interrupt signal")

	// Stop the server using a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		log.WithError(err).Fatal("Error stopping server")
	}

	log.Info("Server gracefully stopped")
}
