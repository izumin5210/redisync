package redisync_test

import (
	"context"
	"sync"
	"testing"

	"github.com/izumin5210/redisync"
)

func TestScoreFilter(t *testing.T) {
	defer cleanupTestRedis()

	m := redisync.NewMonitor(pool)
	filter := redisync.NewScoreFilter(pool, m)

	resultCh := make(chan int)

	go func() {
		defer close(resultCh)
		var wg sync.WaitGroup
		defer wg.Wait()
		for i := 0; i < 100; i++ {
			for j := 0; j < 4; j++ {
				wg.Add(1)
				i := i
				go func() {
					defer wg.Done()
					ok, err := filter.Filter(context.Background(), "foo", i)
					if err != nil {
						t.Errorf("returned %v, want nil", err)
					}
					if ok {
						resultCh <- i
					}
				}()
			}
		}
	}()

	var msgs []int
	for v := range resultCh {
		if len(msgs) > 0 {
			if prev, next := msgs[len(msgs)-1], v; prev >= next {
				t.Errorf("invalid resut: prev %v, next %v", prev, next)
			}
		}
		msgs = append(msgs, v)
	}

	if len(msgs) < 3 {
		t.Errorf("proceeded messages too few: %v", msgs)
	}
	t.Logf("proceeded messages: %v", msgs)
}
