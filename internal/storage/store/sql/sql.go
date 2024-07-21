package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"github.com/e1m0re/grdn/internal/db/migrations"
	"github.com/e1m0re/grdn/internal/models"
)

var (
	// ErrPathNotSpecified is the error returned when the path parameter passed in NewStore is blank
	ErrPathNotSpecified = errors.New("path cannot be empty")
)

// Store that leverages a database.
type Store struct {
	db   *sqlx.DB
	path string
}

// NewStore initializes the database and creates the schema if it doesn't already exist in the path specified.
func NewStore(path string) (*Store, error) {
	if len(path) == 0 {
		return nil, ErrPathNotSpecified
	}
	store := &Store{path: path}
	var err error
	if store.db, err = sqlx.Open("pgx", path); err != nil {
		return nil, err
	}
	if err = store.db.Ping(); err != nil {
		return nil, err
	}
	if err = store.migrate(); err != nil {
		_ = store.db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) migrate() error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", s.path)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&migrations.Content)
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	return db.Close()
}

// Clear removes all data in storage.
func (s *Store) Clear(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM metrics WHERE id > 0")
	return err
}

// Close closes the connection to the storage.
func (s *Store) Close() error {
	return s.db.Close()
}

// GetAllMetrics returns the list of all metrics.
func (s *Store) GetAllMetrics(ctx context.Context) (*models.MetricsList, error) {
	var metrics models.MetricsList
	err := s.db.SelectContext(ctx, &metrics, "SELECT name, type, delta, value FROM metrics")

	return &metrics, err
}

// GetMetric returns an object Metric.
func (s *Store) GetMetric(ctx context.Context, mType models.MetricType, mName string) (*models.Metric, error) {
	var metric models.Metric
	err := s.db.GetContext(ctx, &metric, `SELECT name, type, delta, value FROM metrics WHERE name = $1 AND type = $2`, mName, mType)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &metric, err
	}
}

// Ping checks the connection to the storage.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping()
}

// Restore loads data from a file.
func (s *Store) Restore(ctx context.Context) error {
	return nil
}

// Save saves data to a file.
func (s *Store) Save(ctx context.Context) error {
	return nil
}

// UpdateMetrics performs batch updates of result values in the store.
func (s *Store) UpdateMetrics(ctx context.Context, metrics models.MetricsList) error {
	if len(metrics) == 0 {
		return nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (name, type, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT(name, type) DO UPDATE SET delta = $3, value = $4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
		if err != nil {
			rollbackErr := tx.Rollback()
			return errors.Join(err, rollbackErr)
		}
	}

	err = tx.Commit()

	return err
}
