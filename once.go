package redisync

import (
	"context"
)

type Once struct {
	Config
	pool Pool
}

func NewOnce(pool Pool, opts ...Option) *Once {
	return &Once{
		Config: createConfig(opts),
		pool:   pool,
	}
}

func (o *Once) Do(ctx context.Context, key string, f func(context.Context) error) error {
	conn, err := o.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = TryLock(conn, key, o.LockExpiration)
	if err != nil {
		return err
	}

	return f(ctx)
}
