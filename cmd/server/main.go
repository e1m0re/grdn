package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/e1m0re/grdn/internal/http-server/handler"
	gzipMiddleware "github.com/e1m0re/grdn/internal/http-server/middleware/gzip"
	loggerMiddleware "github.com/e1m0re/grdn/internal/http-server/middleware/logger"
	"github.com/e1m0re/grdn/internal/logger"
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
func main() {
	parameters := config()

	if err := logger.Initialize(parameters.loggerLevel); err != nil {
		fmt.Print(err)
		return
	}

	store := storage.NewMemStorage()
	router := initRouter(handler.NewHandler(store))

	fmt.Println("Running server on ", parameters.serverAddr)
	log.Fatal(http.ListenAndServe(parameters.serverAddr, router))
}
