package router

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/samims/ecommerceGo/configs"
	"github.com/samims/ecommerceGo/handlers"
	"github.com/samims/ecommerceGo/storage"
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
func (lr *LocalRouter) GetRouter() *mux.Router {

	ph := handlers.NewProduct(lr.l)
	r := mux.NewRouter()

	s := storage.NewFileStorage(lr.l, lr.cfg)
	// Initialize the files handler
	files := handlers.NewFiles(s, lr.l, lr.cfg)

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

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

	// fileAPIs
	fileUploadRouter := r.Methods(http.MethodPost).Subrouter()
	fileUploadRouter.HandleFunc("/upload", files.UploadMultipart)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Serve the index.html file
	r.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/upload.html")
	})

	r.HandleFunc("/images/{unixTime}/{fileName}", files.ServeImage)
	return r
}
