package monitoring

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Monitor struct {
	db      *pgxpool.Pool
	workers map[int64]chan struct{}
	mu      sync.Mutex
}

// NewMonitor creates a new Monitor instance.
func NewMonitor(db *pgxpool.Pool) *Monitor {
	return &Monitor{
		db:      db,
		workers: make(map[int64]chan struct{}),
	}
}

// StartServiceMonitor starts monitoring a service.
