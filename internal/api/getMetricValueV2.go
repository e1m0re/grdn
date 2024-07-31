package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/e1m0re/grdn/internal/models"
)

func (h *Handler) getMetricValueV2(response http.ResponseWriter, request *http.Request) {
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

	metric, err := h.services.MetricsManager.GetMetric(request.Context(), reqData.MType, reqData.ID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if metric == nil {
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
