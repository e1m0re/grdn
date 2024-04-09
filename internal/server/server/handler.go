package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/models"
)

func (srv *Server) updateMetric(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mType := chi.URLParam(request, "mType")
	mName := chi.URLParam(request, "mName")
	mValue := chi.URLParam(request, "mValue")

	err := srv.store.UpdateMetricValue(mType, mName, mValue)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (srv *Server) getMainPage(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/html")

	for _, value := range srv.store.GetAllMetrics() {
		_, err := fmt.Fprintf(response, "%s\r\n", value)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	response.WriteHeader(http.StatusOK)
}

func (srv *Server) getMetricValue(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metric, err := srv.store.GetMetric(chi.URLParam(request, "mType"), chi.URLParam(request, "mName"))
	if err != nil {
		http.Error(response, "Not found.", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "text/html")
	_, err = response.Write([]byte(metric.ValueToString()))

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (srv *Server) getMetricValueV2(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var buf bytes.Buffer

	_, err := buf.ReadFrom(request.Body)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}

	var reqData models.Metrics

	if err = json.Unmarshal(buf.Bytes(), &reqData); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := srv.store.GetMetric(reqData.MType, reqData.ID)

	if err != nil {
		http.Error(response, "Not found.", http.StatusNotFound)
		return
	}

	respContent, err := json.Marshal(metric)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	_, err = response.Write(respContent)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (srv *Server) updateMetrics(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data models.Metrics

	var buf bytes.Buffer

	_, err := buf.ReadFrom(request.Body)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	err = srv.store.UpdateMetricValueV2(data)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (srv *Server) checkDBConnection(response http.ResponseWriter, request *http.Request) {
	if err := srv.store.Ping(request.Context()); err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
