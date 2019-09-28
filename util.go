package redisync

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/gomodule/redigo/redis"
)

var (
	ErrLocked   = errors.New("already locked")
	ErrConflict = errors.New("locked by another process")
)

type BackOffFactory interface {
	Create(ctx context.Context) backoff.BackOff
}

type BackOffFactoryFunc func() backoff.BackOff

func (f BackOffFactoryFunc) Create(ctx context.Context) (bo backoff.BackOff) {
	return backoff.WithContext(f(), ctx)
}

func TryLock(conn redis.Conn, key string, expiration time.Duration) (string, error) {
	r, err := secureRandom(24)
	if err != nil {
		return "", err
	}

	v, err := conn.Do("SET", key, r, "EX", expiration.Seconds(), "NX")
	if err != nil {
		return "", err
	}
	if v == nil {
		return "", ErrLocked
	}

	return r, nil
}

func Unlock(conn redis.Conn, key, value string) (err error) {
	_, err = conn.Do("WATCH", key)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			conn.Do("DISCARD")
		}
	}()

	v, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return err
	}
	if v != value {
		return ErrConflict
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("DEL", key)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	return err
}

func secureRandom(b int) (string, error) {
	k := make([]byte, b)
	if _, err := rand.Read(k); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", k), nil
}
