package redisync

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

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
