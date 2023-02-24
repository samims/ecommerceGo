package router

import (
	"net/http"

	"product-images/configs"
	"product-images/handlers"
	localMiddleware "product-images/middlewares"
	"product-images/storage"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
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

	mw := localMiddleware.GzipHandler{}

	r := mux.NewRouter()

	s := storage.NewFileStorage(lr.l, lr.cfg)
	// Initialize the files handler
	files := handlers.NewFiles(s, lr.l, lr.cfg)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	getRouter := r.Methods(http.MethodGet).Subrouter()

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

	imageAccessRouter := r.Methods(http.MethodGet).Subrouter()
	imageAccessRouter.HandleFunc("/images/{unixTime}/{fileName}", files.ServeImage)
	imageAccessRouter.Use(mw.GzipMiddleware)
	return r
}
