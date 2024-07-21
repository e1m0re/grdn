package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/utils"
)

func (h *Handler) updateMetricsList(response http.ResponseWriter, request *http.Request) {
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

	ctx, cancelFunc := context.WithCancel(request.Context())
	defer cancelFunc()
	err = utils.RetryFunc(ctx, func() error {
		return h.services.MetricsManager.UpdateMetrics(ctx, metrics)
	})
	if err != nil {
		slog.Error("update metrics error", slog.String("error", err.Error()))
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}
