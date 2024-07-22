package service

import (
	"github.com/e1m0re/grdn/internal/agent/config"
	"github.com/e1m0re/grdn/internal/service/apiclient"
	"github.com/e1m0re/grdn/internal/service/encryption"
	"github.com/e1m0re/grdn/internal/service/metrics"
	"github.com/e1m0re/grdn/internal/service/monitor"
	"github.com/e1m0re/grdn/internal/service/storage"
	"github.com/e1m0re/grdn/internal/storage/store"
)

// ServerServices is servers DI-container.
type ServerServices struct {
	MetricsManager metrics.Manager
	StorageService storage.Service
}

// NewServerServices is ServerServices constructor.
func NewServerServices(s store.Store) *ServerServices {
	return &ServerServices{
		MetricsManager: metrics.NewMetricsManager(s),
		StorageService: storage.NewService(s),
	}
}

// AgentServices is agents DI-container.
type AgentServices struct {
	APIClient *apiclient.APIClient
	Monitor   monitor.Monitor
	Encryptor encryption.Encryptor
}

// NewAgentServices is AgentServices constructor.
func NewAgentServices(cfg *config.Config) (*AgentServices, error) {
	var encr encryption.Encryptor
	var err error
	if len(cfg.PublicKeyFile) > 0 {
		encr, err = encryption.NewEncryptor(cfg.PublicKeyFile)
		if err != nil {
			return nil, err
		}
	}

	return &AgentServices{
		APIClient: apiclient.NewAPIClient("http://"+cfg.ServerAddr, []byte(cfg.Key)),
		Monitor:   monitor.NewMonitor(),
		Encryptor: encr,
	}, nil
}
