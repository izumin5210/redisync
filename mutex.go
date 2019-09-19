package redisync

import (
	"context"

	"github.com/cenkalti/backoff/v3"
)

type Mutex struct {
	Config
	pool Pool
	key  string
}

func NewMutex(pool Pool, key string, opts ...Option) *Mutex {
	return &Mutex{
		Config: createConfig(opts),
		pool:   pool,
		key:    key,
	}
}

func (m *Mutex) Lock(ctx context.Context) error {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return backoff.Retry(func() error {
		return TryLock(conn, m.key, m.LockExpiration)
	}, m.BackOffFactory.Create(ctx))
}

func (m *Mutex) Unlock(ctx context.Context) error {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return Unlock(conn, m.key)
}
