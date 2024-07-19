package service

import (
	"github.com/e1m0re/grdn/internal/server/service/metrics"
	"github.com/e1m0re/grdn/internal/server/service/storage"
	"github.com/e1m0re/grdn/internal/server/storage/store"
)

// Services is DI-container.
type Services struct {
	MetricsManager metrics.Manager
	StorageService storage.Service
}

// NewServices is Services constructor.
func NewServices(s store.Store) *Services {
	return &Services{
		MetricsManager: metrics.NewMetricsManager(s),
		StorageService: storage.NewService(s),
	}
}
