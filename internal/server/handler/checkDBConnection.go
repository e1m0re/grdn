package handler

import (
	"github.com/e1m0re/grdn/internal/server/storage/store"
	"log/slog"
	"net/http"
)

func (h *Handler) checkDBConnection(response http.ResponseWriter, request *http.Request) {
	err := store.Get().Ping(request.Context())
	if err != nil {
		slog.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
