package redisync

import (
	"context"
)

type Once struct {
	OnceConfig
	pool Pool
}

func NewOnce(pool Pool, opts ...OnceOption) *Once {
	return &Once{
		OnceConfig: createOnceConfig(opts),
		pool:       pool,
	}
}

func (o *Once) Do(ctx context.Context, key string, f func(context.Context) error) error {
	conn, err := o.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	unlockValue, err := TryLock(conn, key, o.Expiration)
	if err != nil {
		return err
	}

	err = f(ctx)
	if err != nil && o.UnlockAfterError {
		_ = Unlock(conn, key, unlockValue)
	}
	return err
}
