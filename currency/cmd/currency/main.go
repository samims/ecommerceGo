package main

import (
	"fmt"
	"net"
	"os"

	config "github.com/samims/ecommerceGO/currency/configs"
	"github.com/sirupsen/logrus"

	"github.com/samims/ecommerceGO/currency/constants"
	"github.com/samims/ecommerceGO/currency/data"
	"github.com/samims/ecommerceGO/currency/logger"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/samims/ecommerceGO/currency/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.NewLogger(logrus.DebugLevel)
	cfg := config.NewViperConfig()

	gs := grpc.NewServer()
	rates, err := data.NewRates(log, cfg)
	if err != nil {
		log.Error("unable to generate rates", "error", err)
		os.Exit(1)
	}

	cs := server.NewCurrency(log, rates)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)
	portStr := cfg.GetString(constants.EnvPort)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", portStr))
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	log.Println("Serving on port", portStr)
	err = gs.Serve(l)
	if err != nil {
		log.Fatal("unable to serve grpc ", err)
	}

}
