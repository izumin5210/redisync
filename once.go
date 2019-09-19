package redisync

import (
	"context"
	"errors"
)

type Once struct {
	Config
	pool Pool
}

var (
	ErrConflict = errors.New("this operation has been proceeded")
)

func NewOnce(pool Pool, opts ...Option) *Once {
	return &Once{
		Config: createConfig(opts),
		pool:   pool,
	}
}

func (o *Once) Run(ctx context.Context, key string, f func(context.Context) error) error {
	conn, err := o.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = TryLock(conn, key, o.LockExpiration)
	if err != nil {
		return err
	}

	return f(ctx)
}
