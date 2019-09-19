package redisync_test

import (
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func TestMain(m *testing.M) {
	pool = &redis.Pool{
		Dial:      func() (redis.Conn, error) { return redis.DialURL(os.Getenv("REDIS_URL")) },
		MaxIdle:   100,
		MaxActive: 100,
		Wait:      true,
	}
	defer pool.Close()

	code := m.Run()
	defer os.Exit(code)
}

func cleanupTestRedis() {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		panic(err)
	}
}
