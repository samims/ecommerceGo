package router

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/samims/ecommerceGo/configs"
	"github.com/samims/ecommerceGo/handlers"
	"github.com/sirupsen/logrus"
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
func (r *LocalRouter) GetRouter() *mux.Router {

	ph := handlers.NewProduct(r.l)
	sm := mux.NewRouter()

	storage := handlers.NewLocalDiskStorage(r.cfg.FileDir, "id")

	// Initialize the files handler
	files := handlers.NewFiles(storage, r.l, r.cfg)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)

	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// fileAPIs
	fileUploadRouter := sm.Methods(http.MethodPost).Subrouter()
	fileUploadRouter.HandleFunc("/upload", files.UploadSingleFile)

	return sm
}
