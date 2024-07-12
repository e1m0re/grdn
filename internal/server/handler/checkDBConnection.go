package handler

import (
	"context"
	"log/slog"
	"net/http"
)

func (h *Handler) checkDBConnection(response http.ResponseWriter, request *http.Request) {
	ctx, cancelFunc := context.WithCancel(request.Context())
	defer cancelFunc()

	err := h.services.StorageService.TestConnection(ctx)
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
