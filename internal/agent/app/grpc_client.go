package app

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/proto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=ClientGRPC
type ClientGRPC interface {
	// Shutdown closes connection.
	Shutdown() error
	// UpdateMetricsList sends data to server.
	UpdateMetricsList(ctx context.Context, metrics models.MetricsList) error
}

type client struct {
	conn          *grpc.ClientConn
	metricsClient proto.MetricsClient
}

// UpdateMetricsList sends data to server.
func (c *client) UpdateMetricsList(ctx context.Context, metrics models.MetricsList) error {
	count := len(metrics)
	if count == 0 {
		return nil
	}

	request := &proto.UpdateMetricsListRequest{
		MetricsData: make([]*proto.MetricData, count),
		Count:       int32(count),
	}

	for idx, metric := range metrics {
		request.MetricsData[idx] = &proto.MetricData{
			Id:    metric.ID,
			Value: metric.ValueToString(),
		}

		switch metric.MType {
		case models.CounterType:
			request.MetricsData[idx].Type = 0
		case models.GaugeType:
			request.MetricsData[idx].Type = 1
		default:
			return fmt.Errorf("unknown metric type (%s)", metric.MType)
		}
	}

	_, err := c.metricsClient.UpdateMetricsList(ctx, request)
	if err != nil {
		msg := "sending metrics data to server"
		if e, ok := status.FromError(err); ok {
			slog.Warn(msg, slog.String("code", e.Code().String()), slog.String("error", e.Message()))
		} else {
			slog.Warn(msg, slog.String("error", e.String()))
		}
	}

	return nil
}

// Shutdown closes connection.
func (c *client) Shutdown() error {
	return c.conn.Close()
}

var _ ClientGRPC = (*client)(nil)

func NewClientGRPC(address string) (ClientGRPC, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &client{
		conn:          connection,
		metricsClient: proto.NewMetricsClient(connection),
	}, nil
}
