package redisync

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func lock(conn redis.Conn, key string, expiration time.Duration) error {
	v, err := conn.Do("SET", key, 1, "EX", expiration.Seconds(), "NX")
	if err != nil {
		return err
	}
	if v == nil {
		return ErrConflict
	}
	return nil
}
