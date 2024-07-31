package api

import (
	"bytes"
	"context"
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

	var metric models.Metric
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancelFunc := context.WithCancel(request.Context())
	defer cancelFunc()
	err = utils.RetryFunc(ctx, func() error {
		return h.services.MetricsManager.UpdateMetric(ctx, metric)
	})
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}
