package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) checkDBConnection(response http.ResponseWriter, request *http.Request) {
	err := h.services.StorageService.PingDB(request.Context())
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
