// Package server provides a struct for HTTP server and methods to run it
// and handle graceful shutdown.

package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"product-images/configs"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

// Server holds an HTTP server instance and router instance
type Server struct {
	Router http.Handler
	Srv    *http.Server
	log    *logrus.Logger
}

// NewServer creates and returns a new instance of Server
func NewServer(handler http.Handler, cfg *configs.Config, l *logrus.Logger) *Server {
	ch := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins(cfg.AllowedHosts))

	return &Server{
		Router: handler,
		log:    l,
		Srv: &http.Server{
			Addr:         cfg.ServerCfg.Addr,
			Handler:      ch(handler),
			IdleTimeout:  cfg.ServerCfg.IdleTimeOut,
			ReadTimeout:  cfg.ServerCfg.ReadTimeOut,
			WriteTimeout: cfg.ServerCfg.WriteTimeOut,
		},
	}
}

// GraceFulShutDown waits for an interrupt signal and gracefully shuts down the server
func (s *Server) GraceFulShutDown(killTime time.Duration) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)
	signal.Notify(stopCh, os.Kill)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), killTime)

	defer cancel()

	s.log.Infoln("Shutting down server...")
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	return s.Srv.ListenAndServe()
}

// Shutdown shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.Srv.Shutdown(ctx)
}
