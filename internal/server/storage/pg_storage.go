package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/e1m0re/grdn/internal/models"
)

type DBStorage struct {
	db *sqlx.DB
}

// NewDBStorage is DNStorage constructor.
func NewDBStorage(dsn string) (*DBStorage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db: db,
	}, nil
}

// Close closes the connection to the storage.
func (s *DBStorage) Close() error {
	return s.db.Close()
}

// DumpStorageToFile saves data to a file.
func (s *DBStorage) DumpStorageToFile() error {
	return nil
}

// GetMetric returns an object Metric.
func (s *DBStorage) GetMetric(ctx context.Context, mType models.MetricType, mName string) (metric *models.Metric, err error) {
	metric = &models.Metric{}
	err = s.db.GetContext(ctx, metric, `SELECT name, type, delta, value FROM metrics WHERE name = $1 AND type = $2`, mName, mType)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrUnknownMetric
	case err != nil:
		return nil, err
	default:
		return
	}
}

// GetMetricsList returns a list of metrics in the format <METRIC>:<VALUE>.
func (s *DBStorage) GetMetricsList(ctx context.Context) ([]string, error) {

	var metrics models.MetricsList
	err := s.db.SelectContext(ctx, &metrics, "SELECT name, type, delta, value FROM metrics")
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, metric := range metrics {
		result = append(result, metric.String())
	}

	return result, nil
}

// LoadStorageFromFile loads data from a file.
func (s *DBStorage) LoadStorageFromFile() error {
	return nil
}

// Ping checks the connection to the storage.
func (s *DBStorage) Ping(ctx context.Context) error {
	return s.db.Ping()
}

// UpdateMetric performs updates to the value of the specified metric in the store.
func (s *DBStorage) UpdateMetric(ctx context.Context, metric models.Metric) error {
	query := `INSERT INTO metrics (name, type, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT(name, type) DO UPDATE SET delta = (CASE WHEN metrics.delta IS NULL THEN NULL ELSE metrics.delta + $3 END), value = $4`
	_, err := s.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value)
	return err
}

// UpdateMetrics performs batch updates of metric values in the store.
func (s *DBStorage) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	if len(metrics) == 0 {
		return nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		err = s.UpdateMetric(ctx, *metric)
		if err != nil {
			rollbackErr := tx.Rollback()
			return errors.Join(err, rollbackErr)
		}
	}

	err = tx.Commit()

	return err
}
