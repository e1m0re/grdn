package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/utils"
)

func (h *Handler) updateMetric(response http.ResponseWriter, request *http.Request) {
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

	ctx, cancelFunc := context.WithCancel(request.Context())
	defer cancelFunc()
	err = utils.RetryFunc(ctx, func() error {
		return h.services.MetricsManager.UpdateMetric(ctx, metric)
	})

	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}
