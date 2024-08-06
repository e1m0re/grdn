package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	metricsgrpc "github.com/e1m0re/grdn/internal/grpc/metrics"
	"github.com/e1m0re/grdn/internal/listeners"
	"github.com/e1m0re/grdn/internal/service"
)

type listener struct {
	gRPCServer *grpc.Server
	config     *Config
}

// Run starts GRPC listener.
func (l *listener) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", l.config.ServerAddr)
	if err != nil {
		return err
	}

	return l.gRPCServer.Serve(listen)
}

// Shutdown stops GRPC listener.
func (l *listener) Shutdown(ctx context.Context) error {
	l.gRPCServer.GracefulStop()
	return nil
}

var _ listeners.Listener = (*listener)(nil)

// NewGRPSListener initiates new instance of GRPC listener.
func NewGRPSListener(cfg *Config, services service.ServerServices) (listeners.Listener, error) {
	listener := &listener{
		gRPCServer: grpc.NewServer(),
		config:     cfg,
	}

	metricsGRPCService := metricsgrpc.NewService(services.MetricsManager)
	metricsGRPCService.Register(listener.gRPCServer)

	return listener, nil
}
