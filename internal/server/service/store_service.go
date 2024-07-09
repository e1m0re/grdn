package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/server/storage"
)

type storageService struct {
	store storage.Storage
}

// NewStorageService is storageService Constructor.
func NewStorageService(store storage.Storage) StorageService {
	return &storageService{
		store: store,
	}
}

// Close closes the connection to the storage.
func (ss *storageService) Close() error {
	return ss.store.Close()
}

// DumpStorageToFile saves data to a file.
func (ss *storageService) DumpStorageToFile() error {
	return ss.store.DumpStorageToFile()
}

// LoadStorageFromFile loads data from a file.
func (ss *storageService) LoadStorageFromFile() error {
	return ss.store.LoadStorageFromFile()
}

// PingDB checks the connection to the storage.
func (ss *storageService) PingDB(ctx context.Context) error {
	return ss.store.Ping(ctx)
}
