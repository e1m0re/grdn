package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/storage"
)

var store = storage.NewMemStorage()

func updateMetricHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mType := chi.URLParam(request, "mType")
	if !storage.IsValidMetricsType(mType) {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := chi.URLParam(request, "mName")
	mValue := chi.URLParam(request, "mValue")

	switch mType {
	case storage.GaugeType:
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateGaugeMetric(mName, value)
	case storage.CounterType:
		value, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateCounterMetric(mName, value)
	}
}

func mainPageHandler(response http.ResponseWriter, _ *http.Request) {
	for _, value := range store.GetAllMetrics() {
		_, err := fmt.Fprintf(response, "%s\r\n", value)
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
	}
}

func getMetricValueHandler(response http.ResponseWriter, request *http.Request) {
	value, err := store.GetMetricValue(chi.URLParam(request, "mType"), chi.URLParam(request, "mName"))
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	}

	_, err = response.Write([]byte(value))
	if err != nil {
		fmt.Printf("%v\r\n", err)
	}
}

func AppRouter() chi.Router {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", mainPageHandler)
		router.Get("/value/{mType}/{mName}", getMetricValueHandler)
		router.Post("/update/{mType}/{mName}/{mValue}", updateMetricHandler)
	})
	return router
}

func main() {
	parseFlags()

	fmt.Println("Running server on", flagRunAddr)
	log.Fatal(http.ListenAndServe(flagRunAddr, AppRouter()))
}
