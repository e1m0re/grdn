package storage

import "time"

// Config is the configuration for storage
type Config struct {
	// Path is the path used by the store to achieve persistence
	// If blank, persistence is disabled.
	Path string

	// Type of store
	Type Type

	// SyncMode toggle of autosave mode
	SyncMode bool

	// Interval of autosave in-memory store
	Interval time.Duration
}
