package srvhandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/storage"
)

type Handler struct {
	store *storage.MemStorage
}

func NewHandler(store *storage.MemStorage) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) UpdateMetric(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mType := chi.URLParam(request, "mType")
	mName := chi.URLParam(request, "mName")
	mValue := chi.URLParam(request, "mValue")

	err := h.store.UpdateMetricValue(mType, mName, mValue)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetMainPage(response http.ResponseWriter, _ *http.Request) {
	for _, value := range h.store.GetAllMetrics() {
		_, err := fmt.Fprintf(response, "%s\r\n", value)
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
	}
}

func (h *Handler) GetMetricValue(response http.ResponseWriter, request *http.Request) {
	value, err := h.store.GetMetricValue(chi.URLParam(request, "mType"), chi.URLParam(request, "mName"))
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	}

	_, err = response.Write([]byte(value))
	if err != nil {
		fmt.Printf("%v\r\n", err)
	}
}
