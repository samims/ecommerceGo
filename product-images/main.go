package main

import (
	"time"

	"product-images/configs"
	"product-images/logger"
	"product-images/router"
	"product-images/server"

	"github.com/sirupsen/logrus"
)

var bindAddress = ":8080"
var logLevel = logrus.DebugLevel
var imagedDIR = "./tmp/images"
var mediaURL = "/images"
var allowedHosts = []string{"http://localhost:8000"}

func main() {

	l := logger.NewLogger(logLevel)

	sCfg := configs.NewServerConf(bindAddress, allowedHosts, 120*time.Second, 15*time.Second, 15*time.Second)
	cfg := configs.NewConfig(allowedHosts, imagedDIR, mediaURL, sCfg)

	r := router.NewLocalRouter(l, cfg)
	routerHandler := r.GetRouter()

	s := server.NewServer(routerHandler, cfg, l)

	go func(s *server.Server, l *logrus.Logger) {
		l.Infoln("Starting the handlers on port ", s.Srv.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}(s, l)

	s.GraceFulShutDown(10 * time.Second)

}
