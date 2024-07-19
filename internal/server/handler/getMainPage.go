package handler

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) getMainPage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")

	metrics, err := h.services.MetricsManager.GetAllMetrics(request.Context())
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, metric := range *metrics {
		_, err := fmt.Fprintf(response, "%s\r\n", metric.String())
		if err != nil {
			slog.Error(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
		}
	}

	response.WriteHeader(http.StatusOK)
}
