package sql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/e1m0re/grdn/internal/models"
)

var (
	delta = int64(100)
	value = float64(100.1)
)

func TestNewStore(t *testing.T) {
	if _, err := NewStore("", ""); !errors.Is(err, ErrDatabaseDriverNotSpecified) {
		t.Error("expected error due to blank driver parameter ")
	}
	if _, err := NewStore("pgx", ""); !errors.Is(err, ErrPathNotSpecified) {
		t.Error("expected error due to blank path parameter ")
	}
}

func TestStore_Close(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type want struct {
		err error
	}
	tests := []struct {
		mock func()
		want want
		name string
	}{
		{
			name: "something wrong",
			want: want{
				err: errors.New("something wrong"),
			},
			mock: func() {
				mock.ExpectClose().WillReturnError(errors.New("something wrong"))
			},
		},
		{
			name: "successfully case",
			mock: func() {
				mock.ExpectClose().WillReturnError(nil)
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err = s.Close()
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_GetAllMetrics(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type args struct {
		ctx context.Context
	}
	type want struct {
		err     error
		metrics *models.MetricsList
	}
	tests := []struct {
		mock func()
		want want
		args args
		name string
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:     errors.New("something wrong"),
				metrics: nil,
			},
			mock: func() {
				mock.
					ExpectQuery("SELECT name, type, delta, value FROM metrics").
					WillReturnError(errors.New("something wrong"))
			},
		},
		{
			name: "empty list",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err:     nil,
				metrics: &models.MetricsList{},
			},
			mock: func() {
				rows := sqlxmock.NewRows(make([]string, 0))
				mock.
					ExpectQuery("SELECT name, type, delta, value FROM metrics").
					WillReturnRows(rows)
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
				metrics: &models.MetricsList{
					{
						ID:    "metric 1",
						MType: models.CounterType,
						Delta: &delta,
						Value: nil,
					},
					{
						ID:    "metric 2",
						MType: models.GaugeType,
						Delta: nil,
						Value: &value,
					},
				},
			},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"name", "type", "delta", "value"}).
					AddRow("metric 1", "counter", 100, nil).
					AddRow("metric 2", "gauge", nil, 100.1)
				mock.
					ExpectQuery("SELECT name, type, delta, value FROM metrics").
					WillReturnRows(rows)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err := s.GetAllMetrics(test.args.ctx)
			require.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.metrics, got)
		})
	}
}

func TestStore_GetMetric(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type args struct {
		ctx   context.Context
		mType models.MetricType
		mName string
	}
	type want struct {
		err    error
		metric *models.Metric
	}
	tests := []struct {
		mock func()
		args args
		want want
		name string
	}{
		{
			name: "something wrong",
			args: args{
				ctx:   context.Background(),
				mType: models.CounterType,
				mName: "metric 1",
			},
			want: want{
				err:    errors.New("something wrong"),
				metric: nil,
			},
			mock: func() {
				mock.
					ExpectQuery("^SELECT name, type, delta, value FROM metrics WHERE name = \\$1 AND type = \\$2$").
					WillReturnError(errors.New("something wrong"))
			},
		},
		{
			name: "metric not found",
			args: args{
				ctx:   context.Background(),
				mType: models.CounterType,
				mName: "metric 1",
			},
			want: want{
				err:    nil,
				metric: nil,
			},
			mock: func() {
				mock.
					ExpectQuery("^SELECT name, type, delta, value FROM metrics WHERE name = \\$1 AND type = \\$2$").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx:   context.Background(),
				mType: models.CounterType,
				mName: "metric 1",
			},
			want: want{
				err: nil,
				metric: &models.Metric{
					Value: &value,
					Delta: nil,
					MType: models.GaugeType,
					ID:    "metric 1",
				},
			},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"name", "type", "delta", "value"}).
					AddRow("metric 1", "gauge", nil, 100.1)
				mock.
					ExpectQuery("^SELECT name, type, delta, value FROM metrics WHERE name = \\$1 AND type = \\$2$").
					WillReturnRows(rows)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err := s.GetMetric(test.args.ctx, test.args.mType, test.args.mName)
			require.Equal(t, test.want.err, err)
			if test.want.metric != nil {
				assert.Equal(t, test.want.metric.ID, got.ID)
				assert.Equal(t, test.want.metric.MType, got.MType)
				assert.Equal(t, *test.want.metric.Value, *got.Value)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestStore_Clear(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mock func()
		args args
		want want
		name string
	}{
		{
			name: "something wrong",
			mock: func() {
				mock.
					ExpectExec("DELETE FROM metrics WHERE id > 0").
					WillReturnError(errors.New("something wrong"))
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "successfully case",
			mock: func() {
				mock.
					ExpectExec("DELETE FROM metrics WHERE id > 0").
					WillReturnResult(sqlxmock.NewResult(0, 0))
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := s.Clear(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_Ping(t *testing.T) {
	db, mock, err := sqlxmock.Newx(sqlxmock.MonitorPingsOption(true))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mock func()
		args args
		want want
		name string
	}{
		{
			name: "something wrong",
			mock: func() {
				mock.
					ExpectPing().
					WillReturnError(errors.New("something wrong"))
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "successfully case",
			mock: func() {
				mock.
					ExpectPing().
					WillReturnError(nil)
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := s.Ping(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestStore_Restore(t *testing.T) {
	s := &Store{}
	err := s.Restore(context.Background())
	assert.Nil(t, err)
}

func TestStore_Save(t *testing.T) {
	s := &Store{}
	err := s.Save(context.Background())
	assert.Nil(t, err)
}

func TestStore_UpdateMetrics(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := Store{db: db}

	type args struct {
		ctx     context.Context
		metrics models.MetricsList
	}
	type want struct {
		err error
	}
	tests := []struct {
		mock func()
		name string
		want want
		args args
	}{
		{
			name: "empty list in args",
			mock: func() {},
			args: args{
				ctx:     context.Background(),
				metrics: make(models.MetricsList, 0),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "BeginTxx failed",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(errors.New("something wrong"))
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{},
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "PrepareContext failed",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(nil)
				mock.
					ExpectPrepare("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(errors.New("something wrong"))
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{},
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "ExecContext failed case",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(nil)
				mock.
					ExpectPrepare("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(nil)
				mock.
					ExpectExec("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(errors.New("something wrong"))
				mock.
					ExpectRollback().
					WillReturnError(nil)
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						ID:    "metric 1",
						MType: models.CounterType,
						Delta: &delta,
					},
				},
			},
			want: want{
				err: errors.Join(errors.New("something wrong"), nil),
			},
		},
		{
			name: "ExecContext and Rollback failed case",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(nil)
				mock.
					ExpectPrepare("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(nil)
				mock.
					ExpectExec("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(errors.New("something wrong"))
				mock.
					ExpectRollback().
					WillReturnError(errors.New("something wrong in rollback"))
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						ID:    "metric 1",
						MType: models.CounterType,
						Delta: &delta,
					},
				},
			},
			want: want{
				err: errors.Join(errors.New("something wrong"), errors.New("something wrong in rollback")),
			},
		},
		{
			name: "Commit failed case",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(nil)
				mock.
					ExpectPrepare("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(nil)
				mock.
					ExpectExec("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.
					ExpectCommit().
					WillReturnError(errors.New("something wrong"))
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						ID:    "metric 1",
						MType: models.CounterType,
						Delta: &delta,
					},
				},
			},
			want: want{
				err: errors.New("something wrong"),
			},
		},
		{
			name: "successfully case",
			mock: func() {
				mock.
					ExpectBegin().
					WillReturnError(nil)
				mock.
					ExpectPrepare("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnError(nil)
				mock.
					ExpectExec("^INSERT INTO metrics \\(name, type, delta, value\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT\\(name, type\\) DO UPDATE SET delta = \\$3, value = \\$4$").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.
					ExpectCommit().
					WillReturnError(nil)
			},
			args: args{
				ctx: context.Background(),
				metrics: models.MetricsList{
					&models.Metric{
						ID:    "metric 1",
						MType: models.CounterType,
						Delta: &delta,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := s.UpdateMetrics(test.args.ctx, test.args.metrics)
			require.Equal(t, test.want.err, err)
		})
	}
}
