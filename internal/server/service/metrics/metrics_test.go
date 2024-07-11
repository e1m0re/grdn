package metrics

import (
	"context"
	"github.com/e1m0re/grdn/internal/models"
	"github.com/e1m0re/grdn/internal/server/storage/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_metricsManager_GetAllMetrics(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		list   *models.MetricsList
		errMsg string
	}
	tests := []struct {
		name      string
		args      args
		mockStore func(ctx context.Context)
		want      want
	}{
		{
			name: "Empty list",
			args: args{
				ctx: context.Background(),
			},
			mockStore: func(ctx context.Context) {
				_ = store.Get().Clear(ctx)
			},
			want: want{
				list:   &models.MetricsList{},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockStore(test.args.ctx)
			mm := &metricsManager{}
			got, err := mm.GetAllMetrics(test.args.ctx)
			if len(test.want.errMsg) > 0 {
				require.Errorf(t, err, test.want.errMsg)
			}
			require.Equal(t, test.want.list, got)
		})
	}
}
