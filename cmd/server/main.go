package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/e1m0re/grdn/internal/storage"
)

var store = storage.NewMemStorage()

func isValidMetricName(mType storage.MetricsType, value string) bool {
	if mType == storage.GuageType {
		return storage.IsValidGuageName(value)
	}
	if mType == storage.CounterType {
		return storage.IsValidCounterName(value)
	}
	return false
}

func updateMetricHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// check type
	pathParams := strings.Split(request.RequestURI, "/")
	if len(pathParams) < 3 || !storage.IsValidMetricsType(pathParams[2]) {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// check name
	if len(pathParams) < 4 || !isValidMetricName(pathParams[2], pathParams[3]) {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// check value
	if len(pathParams) < 5 {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	switch pathParams[2] {
	case storage.GuageType:
		value, err := strconv.ParseFloat(pathParams[4], 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		store.UpdateGuageMetric(pathParams[3], value)
	case storage.CounterType:
		value, err := strconv.ParseInt(pathParams[4], 10, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		store.UpdateCounterMetric(pathParams[3], value)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
