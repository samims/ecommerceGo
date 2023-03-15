package main

import (
	"time"

	"product-api/configs"
	"product-api/constants"
	"product-api/logger"
	"product-api/router"
	"product-api/server"

	"github.com/gorilla/mux"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	envs := configs.NewEnv()

	bindAddress := envs.GetString(constants.BindAddress)
	logLevel := envs.GetString(constants.LogLevel)
	imageDir := envs.GetString(constants.ImageDir)
	mediaURL := envs.GetString(constants.MediaURL)
	allowedHosts := envs.GetStringSlice(constants.AllowedHosts)
	currencyServerBase := envs.GetString(constants.CurrencyServerBase)

	// Initialize the logger.
	l := initLogger(logLevel)

	// app cfg
	appCfg := createAppConfig(allowedHosts, imageDir, mediaURL, currencyServerBase)
	// Create the server configuration.
	serverCfg := createServerConfig(bindAddress)

	// Create the application configuration.
	cfg := configs.NewConfig(serverCfg, appCfg).(configs.Config)

	cc, conn, _ := getCurrencyGrpcClient(cfg)

	defer conn.Close()

	// Create the router.
	r := createRouter(l, &cfg, cc)

	// Create the server.
	routerObj := r.GetRouter()
	s := createServer(routerObj, cfg, l)

	// Start the server.
	startServer(s, l)
}

func createAppConfig(allowedHosts []string, imageDIR, mediaURL, currencyServerBase string) configs.AppConfig {
	return configs.NewAppConfig(allowedHosts, imageDIR, mediaURL, currencyServerBase)

}

func initLogger(logLevel string) *logrus.Logger {
	lLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	return logger.NewLogger(lLevel)
}

func createServerConfig(bindAddress string) configs.ServerConfig {
	return configs.NewServerConf(bindAddress, 10*time.Second, 15*time.Second, 15*time.Second)
}

func getCurrencyGrpcClient(cfg configs.Config) (protos.CurrencyClient, *grpc.ClientConn, error) {
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(cfg.AppConfig().GetCurrencyServerBase(), grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
		//return nil, nil, err
	}
	currencyGrpcClient := protos.NewCurrencyClient(conn)
	return currencyGrpcClient, conn, nil
}

func createRouter(l *logrus.Logger, cfg *configs.Config, cc protos.CurrencyClient) *router.Router {
	r := router.NewRouter(l, cfg, cc)
	return r
}

func createServer(r *mux.Router, cfg configs.Config, l *logrus.Logger) *server.Server {
	return server.NewServer(r, cfg, l)
}

func startServer(s *server.Server, l *logrus.Logger) {
	go func(s *server.Server, l *logrus.Logger) {
		l.Infoln("Starting the server on port ", s.Srv.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}(s, l)

	s.GraceFulShutDown(10 * time.Second)
}
