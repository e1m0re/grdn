package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	data map[string]int
}

func (store *MemStorage) UpdateMetric(mType string, name string, value string) {
	i, _ := strconv.Atoi(value)
	switch mType {
	case "gauge", "float64":
		store.data[name] = i
	case "counter", "int64":
		store.data[name] += i
	}
}

var storage = MemStorage{data: map[string]int{}}

func isValidMetricsType(value string) bool {
	switch value {
	case "gauge", "float64", "counter", "int64":
		return true
	default:
		return false
	}
}

func isValidValue(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func updateMetricHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// check type
	pathParams := strings.Split(request.RequestURI, "/")
	if len(pathParams) < 3 || !isValidMetricsType(pathParams[2]) {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// check name
	if len(pathParams) < 4 || len(pathParams[3]) == 0 {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// check value
	if len(pathParams) < 5 || !isValidValue(pathParams[4]) {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	storage.UpdateMetric(pathParams[2], pathParams[3], pathParams[4])

	return
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
