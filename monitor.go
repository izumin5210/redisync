package redisync

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/gomodule/redigo/redis"
)

var (
	synchronizeTimeout = 120 * time.Second
)

func NewMonitor(pool *redis.Pool, bo backoff.BackOff) *Monitor {
	return &Monitor{
		pool: pool,
		bo:   bo,
	}
}

type Monitor struct {
	pool Pool
	bo   backoff.BackOff
}

func (m *Monitor) Synchronize(ctx context.Context, key string, do func(context.Context) error) (err error) {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = backoff.Retry(func() error {
		return TryLock(conn, key, synchronizeTimeout)
	}, m.bo)
	if err != nil {
		return err
	}
	defer func() {
		dErr := Unlock(conn, key)
		if dErr != nil && err == nil {
			err = dErr
		}
	}()

	return do(ctx)
}
