package service

import (
	"context"

	"github.com/e1m0re/grdn/internal/server/storage"
)

type storageService struct {
	store storage.Storage
}

func NewStorageService(store storage.Storage) StorageService {
	return &storageService{
		store: store,
	}
}

func (ss *storageService) PingDB(ctx context.Context) error {
	return ss.store.Ping(ctx)
}

func (ss *storageService) DumpStorageToFile() error {
	return ss.store.DumpStorageToFile()
}

func (ss *storageService) LoadStorageFromFile() error {
	return ss.store.LoadStorageFromFile()
}

func (ss *storageService) Close() error {
	return ss.store.Close()
}
