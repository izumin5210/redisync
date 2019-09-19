package redisync

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/gomodule/redigo/redis"
)

type BackOffFactory interface {
	Create(ctx context.Context) backoff.BackOff
}

type BackOffFactoryFunc func() backoff.BackOff

func (f BackOffFactoryFunc) Create(ctx context.Context) (bo backoff.BackOff) {
	return backoff.WithContext(f(), ctx)
}

func TryLock(conn redis.Conn, key string, expiration time.Duration) error {
	v, err := conn.Do("SET", key, 1, "EX", expiration.Seconds(), "NX")
	if err != nil {
		return err
	}
	if v == nil {
		return ErrConflict
	}
	return nil
}

func Unlock(conn redis.Conn, key string) error {
	_, err := conn.Do("DEL", key)
	return err
}
