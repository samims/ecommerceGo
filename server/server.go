package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/samims/ecommerceGo/configs"
)

type Server struct {
	Router http.Handler
	Srv    *http.Server
}

func NewServer(handler http.Handler, cfg *configs.Config) *Server {
	ch := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins(cfg.AllowedHosts))

	return &Server{
		Router: handler,
		Srv: &http.Server{
			Addr:         cfg.ServerCfg.Addr,
			Handler:      ch(handler),
			IdleTimeout:  cfg.ServerCfg.IdleTimeOut,
			ReadTimeout:  cfg.ServerCfg.ReadTimeOut,
			WriteTimeout: cfg.ServerCfg.WriteTimeOut,
		},
	}
}

func (s *Server) GraceFulShutDown(killTime time.Duration) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)
	signal.Notify(stopCh, os.Kill)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), killTime)

	defer cancel()

	log.Printf("Shutting down server...")
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

}

//func NewServer(Router http.Handler, conf *ServerConf) *Server {
//	return &Server{
//		Router: Router,
//		Srv: &http.Server{
//			Addr:         conf.Addr,
//			Handler:      Router,
//			IdleTimeout:  conf.IdleTimeOut,
//			ReadTimeout:  conf.ReadTimeOut,
//			WriteTimeout: conf.WriteTimeOut,
//		},
//	}
//}

func (s *Server) ListenAndServe() error {
	return s.Srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Srv.Shutdown(ctx)
}
