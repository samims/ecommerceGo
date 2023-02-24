package main

import (
	"fmt"
	"time"

	"product-api/configs"
	"product-api/logger"
	"product-api/router"
	"product-api/server"

	"github.com/sirupsen/logrus"
)

var bindAddress = ":9090"
var logLevel = logrus.DebugLevel
var imagedDIR = "./tmp/images"
var mediaURL = "/images"
var allowedHosts = []string{"http://localhost:8000"}

func main() {

	l := logger.NewLogger(logLevel)

	sCfg := configs.NewServerConf(":9090", allowedHosts, 120*time.Second, 15*time.Second, 15*time.Second)
	cfg := configs.NewConfig(allowedHosts, imagedDIR, mediaURL, sCfg)

	r := router.NewLocalRouter(l, cfg)
	routerHandler := r.GetRouter()

	s := server.NewServer(routerHandler, cfg)

	go func(s *server.Server) {
		fmt.Println("Starting the server on port ", s.Srv.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}(s)

	s.GraceFulShutDown(10 * time.Second)

}
