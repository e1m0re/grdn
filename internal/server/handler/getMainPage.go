package handler

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) getMainPage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")

	sml, err := h.services.MetricsManager.GetSimpleMetricsList(request.Context())
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, value := range sml {
		_, err := fmt.Fprintf(response, "%s\r\n", value)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	response.WriteHeader(http.StatusOK)
}
