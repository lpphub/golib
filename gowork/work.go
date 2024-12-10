package gowork

import (
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

const (
	// DefaultAntsPoolSize sets up the capacity of worker pool, 256 * 1024.
	DefaultAntsPoolSize = 1 << 18

	// ExpiryDuration is the interval time to clean up those expired workers.
	ExpiryDuration = 10 * time.Second

	// Nonblocking decides what to do when submitting a new task to a full worker pool: waiting for a available worker
	// or returning nil directly.
	Nonblocking = true
)

func init() {
	// It releases the default pool from ants.
	ants.Release()
}

// Pool is the alias of ants.Pool.
type Pool = ants.Pool

var (
	defaultAntsPool *Pool
	once            sync.Once
)

// Default instantiates a non-blocking *WorkerPool with the capacity of DefaultAntsPoolSize.
func Default() *Pool {
	once.Do(func() {
		options := defaultOpt()
		defaultAntsPool, _ = ants.NewPool(DefaultAntsPoolSize, ants.WithOptions(options))
	})
	return defaultAntsPool
}

func NewPool(size int) *Pool {
	if size <= 0 || size > DefaultAntsPoolSize {
		return defaultAntsPool
	}
	options := defaultOpt()
	pool, _ := ants.NewPool(size, ants.WithOptions(options))
	return pool
}

func defaultOpt() ants.Options {
	return ants.Options{
		ExpiryDuration: ExpiryDuration,
		Nonblocking:    Nonblocking,
		Logger:         logger.Log(),
		PanicHandler: func(i interface{}) {
			logger.Log().Error().Msgf("goroutine pool panic: %v", i)
		},
	}
}
