package memory

import (
	"context"
	"github.com/e1m0re/grdn/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore_Clear(t *testing.T) {
	type fields struct {
		metrics map[string]models.Metric
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		err          error
		metricsCount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Successfully case",
			fields: fields{
				metrics: map[string]models.Metric{"metric1": {ID: "metric1"}},
			},
			args: args{ctx: nil},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Store{
				metrics: test.fields.metrics,
			}

			err := s.Clear(test.args.ctx)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.metricsCount, len(s.metrics))
		})
	}
}
