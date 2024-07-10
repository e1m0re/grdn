package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/server/storage"
)

func (h *Handler) getMetricValue(response http.ResponseWriter, request *http.Request) {
	metric, err := h.services.MetricsManager.GetMetric(request.Context(), chi.URLParam(request, "mType"), chi.URLParam(request, "mName"))
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
