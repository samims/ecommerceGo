package router

import (
	"net/http"

	"product-api/configs"
	"product-api/data"
	"product-api/handlers"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	protos "github.com/samims/ecommerceGO/currency/protos/currency"
	"github.com/sirupsen/logrus"
)

type Router struct {
	router *mux.Router
}

func NewRouter(logger *logrus.Logger, cfg *configs.Config, cc protos.CurrencyClient) *Router {
	logger.Infof("Router is being initialized with config: %+v", *cfg)

	router := mux.NewRouter()

	pdb := data.NewProductsDB(cc, logger)
	ph := handlers.NewProduct(logger, pdb)

	registerRoutes(router, ph)

	return &Router{
		router: router,
	}
}

// GetRouter returns the pointer to the router instance.
// The router is ready to use for the HTTP server.
func (r *Router) GetRouter() *mux.Router {
	return r.router
}

func registerRoutes(router *mux.Router, ph *handlers.Products) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetByID)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	registerDocs(router)
	registerStatic(router)
}

func registerDocs(router *mux.Router) {
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
}

func registerStatic(router *mux.Router) {
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/upload.html")
	})
}
