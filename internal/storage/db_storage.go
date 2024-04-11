package storage

import (
	"context"
	"database/sql"
	"github.com/e1m0re/grdn/internal/models"
	"log/slog"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, dsn string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db: db,
	}, nil
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}

func (s *DBStorage) DumpStorageToFile() error {
	return nil
}
func (s *DBStorage) GetAllMetrics(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query("SELECT name, type FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		row := make([]string, 2)
		err := rows.Scan(&row[0], &row[1])
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		result = append(result, "%s: %s", row[0], row[1])
	}

	return result, nil
}
func (s *DBStorage) GetMetric(mType models.MetricsType, mName string) (metric *models.Metrics, err error) {
	metric = &models.Metrics{}
	return metric, nil
}
func (s *DBStorage) LoadStorageFromFile() error {
	return nil
}
func (s *DBStorage) Ping(ctx context.Context) error {
	return s.db.Ping()
}
func (s *DBStorage) UpdateCounterMetric(name CounterName, value CounterDateType) {
	return
}
func (s *DBStorage) UpdateGaugeMetric(name GaugeName, value GaugeDateType) {
	return
}
func (s *DBStorage) UpdateMetricValue(mType models.MetricsType, mName string, mValue string) error {
	return nil
}
func (s *DBStorage) UpdateMetricValueV2(data models.Metrics) error {
	return nil
}
