package redisync

import (
	"context"

	"github.com/cenkalti/backoff/v3"
)

func NewMonitor(pool Pool, opts ...Option) *Monitor {
	return &Monitor{
		Config: createConfig(opts),
		pool:   pool,
	}
}

type Monitor struct {
	Config
	pool Pool
}

func (m *Monitor) Synchronize(ctx context.Context, key string, do func(context.Context) error) (err error) {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = backoff.Retry(func() error {
		return TryLock(conn, key, m.LockExpiration)
	}, m.BackOffFactory.Create(ctx))
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
