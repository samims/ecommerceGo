package router

import (
	"net/http"

	"product-api/configs"
	"product-api/data"
	"product-api/handlers"

	//protos "currency/protos/currency"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LocalRouter struct {
	l   *logrus.Logger
	cfg *configs.Config
}

func NewLocalRouter(l *logrus.Logger, cfg *configs.Config) *LocalRouter {
	l.Printf("Router is being initialized with config %+v\n", *cfg)
	r := &LocalRouter{
		l:   l,
		cfg: cfg,
	}
	return r
}

// GetRouter returns a new instance of mux.Router with all the routes and middlewares registered for the local router.
// It returns the pointer to the router instance.
// The router is ready to use for the HTTP server.
func (lr *LocalRouter) GetRouter() *mux.Router {

	r := mux.NewRouter()

	cc := lr.getCurrencyGrpcClient()
	pdb := data.NewProductsDB(cc, lr.l)
	ph := handlers.NewProduct(lr.l, pdb)

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetByID)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)

	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/upload.html")
	})

	return r
}

func (lr *LocalRouter) getCurrencyGrpcClient() protos.CurrencyClient {
	// TODO: it's insecure make it secure
	//creds := credentials.NewTLS(nil)
	//conn, err := grpc.Dial(lr.cfg.CurrencyServerBase, grpc.WithInsecure())
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(lr.cfg.CurrencyServerBase, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	currencyGrpcClient := protos.NewCurrencyClient(conn)
	return currencyGrpcClient
}
