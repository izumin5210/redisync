package redisync

import (
	"context"
	"errors"
	"time"
)

type Once struct {
	pool Pool
}

var (
	ErrConflict    = errors.New("this operation has been proceeded")
	onceExpiration = 3 * 24 * time.Hour
)

func NewOnce(pool Pool) *Once {
	return &Once{
		pool: pool,
	}
}

func (o *Once) Run(ctx context.Context, key string, f func(context.Context) error) error {
	conn, err := o.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = TryLock(conn, key, onceExpiration)
	if err != nil {
		return err
	}

	return f(ctx)
}
