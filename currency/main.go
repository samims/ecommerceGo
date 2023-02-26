package main

import (
	"net"
	"os"

	"currency/logger"
	protos "currency/protos/currency"
	"currency/server"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.NewLogger(logrus.DebugLevel)

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	log.Println("Serving on port 9092...")
	gs.Serve(l)

}
