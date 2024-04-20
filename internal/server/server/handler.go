package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e1m0re/grdn/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/models"
)

func (srv *Server) getMainPage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")

	metrics, err := srv.store.GetAllMetrics(request.Context())
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, value := range metrics {
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

	metric, err := srv.store.GetMetric(request.Context(), chi.URLParam(request, "mType"), chi.URLParam(request, "mName"))
	if err != nil {
		if errors.Is(err, storage.ErrUnknownMetric) {
			http.Error(response, "Not found.", http.StatusNotFound)
			return
		}

		slog.Error(err.Error())
		http.Error(response, err.Error(), http.StatusInternalServerError)
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
		http.Error(response, fmt.Sprintf("error reading body: %s", err), http.StatusBadRequest)
	}

	var reqData models.Metric
	if err = json.Unmarshal(buf.Bytes(), &reqData); err != nil {
		http.Error(response, fmt.Sprintf("error parsing body: %s", err), http.StatusBadRequest)
		return
	}

	metric, err := srv.store.GetMetric(request.Context(), reqData.MType, reqData.ID)
	if err != nil {
		if errors.Is(err, storage.ErrUnknownMetric) {
			http.Error(response, "Not found.", http.StatusNotFound)
			return
		}

		slog.Error(err.Error())
		http.Error(response, err.Error(), http.StatusInternalServerError)
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

func (srv *Server) updateMetricV1(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metric := models.Metric{
		ID:    chi.URLParam(request, "mName"),
		MType: chi.URLParam(request, "mType"),
	}
	err := metric.ValueFromString(chi.URLParam(request, "mValue"))
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = retryFunc(request.Context(), func() error {
		return srv.store.UpdateMetric(request.Context(), metric)
	})
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (srv *Server) updateMetricV2(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	var data models.Metric
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	err = retryFunc(request.Context(), func() error {
		return srv.store.UpdateMetric(request.Context(), data)
	})
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (srv *Server) updateMetrics(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	var metrics models.MetricsList
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	err = retryFunc(request.Context(), func() error {
		return srv.store.UpdateMetrics(request.Context(), metrics)
	})
	if err != nil {
		slog.Error("update metrics error", slog.String("error", err.Error()))
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
