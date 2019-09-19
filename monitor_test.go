package redisync_test

import (
	"context"
	"sync"
	"testing"

	"github.com/izumin5210/redisync"
)

func TestMonitor(t *testing.T) {
	defer cleanupTestRedis()

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
