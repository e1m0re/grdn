package metrics

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/proto"
	"github.com/e1m0re/grdn/internal/service/metrics"
	"github.com/e1m0re/grdn/internal/utils"
)

type serverAPI struct {
	metricsManager metrics.Manager
}

func (s *serverAPI) UpdateMetricsList(ctx context.Context, in *proto.UpdateMetricsListRequest) (*proto.UpdateMetricsListResponse, error) {
	var response proto.UpdateMetricsListResponse

	metricsData := make(models.MetricsList, in.Count)
	for idx, metricData := range in.MetricsData {
		m := models.Metric{ID: metricData.Id}

		switch metricData.Type {
		case 0:
			m.MType = models.CounterType
		case 1:
			m.MType = models.GaugeType
		default:
			return nil, status.Error(codes.InvalidArgument, "unknown metric type")
		}

		if err := m.ValueFromString(metricData.Value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid metric value — %s", err.Error())
		}

		metricsData[idx] = &m
	}

	err := utils.RetryFunc(ctx, func() error {
		return s.metricsManager.UpdateMetrics(ctx, metricsData)
	})
	if err != nil {
		slog.Error("update metrics data failed", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &response, nil
}
