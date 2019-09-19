package redisync_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/izumin5210/redisync"
)

func TestMonitor(t *testing.T) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.DialURL(os.Getenv("REDIS_URL")) },
	}
	defer pool.Close()

	defer func() {
		conn := pool.Get()
		defer conn.Close()
		conn.Do("FLUSHALL")
	}()

	m := redisync.NewMonitor(pool)

	var (
		cnt int
		wg  sync.WaitGroup
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			err := m.Synchronize(context.Background(), "foo", func(context.Context) error {
				cnt += i
				return nil
			})

			if err != nil {
				t.Errorf("returned an error: %v", err)
			}
		}()
	}

	wg.Wait()

	if got, want := cnt, 4950; got != want {
		t.Errorf("count is %d, want %d", got, want)
	}
}
