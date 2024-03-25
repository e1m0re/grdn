package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/e1m0re/grdn/internal/http-server/handler"
	gzipMiddleware "github.com/e1m0re/grdn/internal/http-server/middleware/gzip"
	loggerMiddleware "github.com/e1m0re/grdn/internal/http-server/middleware/logger"
	"github.com/e1m0re/grdn/internal/storage"
)

func initRouter(handler *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(loggerMiddleware.Middleware)
	router.Use(gzipMiddleware.Middleware)
	router.Use(middleware.Compress(5, "text/html", "application/json"))

	router.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMainPage)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handler.GetMetricValueV2)
			r.Get("/{mType}/{mName}", handler.GetMetricValue)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handler.UpdateMetrics)
			r.Post("/{mType}/{mName}/{mValue}", handler.UpdateMetric)
		})
	})

	return router
}

func NewServer(addr string, store *storage.MemStorage) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: initRouter(handler.NewHandler(store)),
	}

	return server
}
