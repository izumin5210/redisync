package redisync

import (
	"context"

	"github.com/cenkalti/backoff/v3"
)

type Mutex struct {
	Config
	pool      Pool
	key       string
	unlockKey string
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

	var unlock string

	err = backoff.Retry(func() (err error) {
		unlock, err = TryLock(conn, m.key, m.LockExpiration)
		return err
	}, m.BackOffFactory.Create(ctx))

	m.unlockKey = unlock

	return err
}

func (m *Mutex) Unlock(ctx context.Context) error {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return Unlock(conn, m.key, m.unlockKey)
}
