package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/srvhandler"
	"github.com/e1m0re/grdn/internal/storage"
)

func main() {
	parseFlags()

	store := storage.NewMemStorage()
	handler := srvhandler.Handler{}
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", handler.GetMainPage(store))
		router.Get("/value/{mType}/{mName}", handler.GetMetricValue(store))
		router.Post("/update/{mType}/{mName}/{mValue}", handler.UpdateMetric(store))
	})

	fmt.Println("Running server on", flagRunAddr)
	log.Fatal(http.ListenAndServe(flagRunAddr, router))
}
