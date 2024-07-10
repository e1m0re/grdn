package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/utils"
)

func (h *Handler) updateMetricV2(response http.ResponseWriter, request *http.Request) {
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

	err = utils.RetryFunc(request.Context(), func() error {
		return h.services.MetricsManager.UpdateMetric(request.Context(), data)
	})
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}
