package service

import (
	"github.com/e1m0re/grdn/internal/server/service/metrics"
)

// Services is DI-container.
type Services struct {
	MetricsManager metrics.Manager
}

// NewServices is Services constructor.
func NewServices() *Services {
	return &Services{
		MetricsManager: metrics.NewMetricsManager(),
	}
}
