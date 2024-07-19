package storage

import (
	"context"

	"github.com/e1m0re/grdn/internal/server/storage/store"
)

// Service is the interface that contains all operations for storage.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Service
type Service interface {
	// Clear removes all data in storage.
	Clear(ctx context.Context) error
	// Close closes the connection to the storage.
	Close() error
	// Restore loads data from a file.
	Restore(ctx context.Context) error
	// Save saves data to a file.
	Save(ctx context.Context) error
	// TestConnection checks connection to storage.
	TestConnection(ctx context.Context) error
}

type service struct {
	store.Store
}

// NewService instantiates new Service.
func NewService(s store.Store) Service {
	return &service{
		s,
	}
}

// Clear removes all data in storage.
func (s *service) Clear(ctx context.Context) error {
	return s.Store.Clear(ctx)
}

// Close closes the connection to the storage.
func (s *service) Close() error {
	return s.Store.Close()
}

// Restore loads data from a file.
func (s *service) Restore(ctx context.Context) error {
	return s.Store.Restore(ctx)
}

// Save saves data to a file.
func (s *service) Save(ctx context.Context) error {
	return s.Store.Save(ctx)
}

// TestConnection checks connection to storage.
func (s *service) TestConnection(ctx context.Context) error {
	return s.Store.Ping(ctx)
}
