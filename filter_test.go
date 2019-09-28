package redisync_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/izumin5210/redisync"
)

func TestScoreFilter(t *testing.T) {
	defer cleanupTestRedis()

	filter := redisync.NewScoreFilter(pool)

	type Message struct {
		Process int
		Score   int
		Time    time.Time
	}

	resultCh := make(chan Message)

	go func() {
		defer close(resultCh)
		var wg sync.WaitGroup
		defer wg.Wait()
		for i := 0; i < 4; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					time.Sleep(50 * time.Microsecond)
					msg := Message{
						Process: i,
						Score:   j,
						Time:    time.Now(),
					}
					ok, err := filter.Filter(context.Background(), "foo", msg.Score)
					if err != nil {
						t.Errorf("returned %v, want nil", err)
					}
					if ok {
						resultCh <- msg
					}
				}
			}()
		}
	}()

	var msgs []Message
	for msg := range resultCh {
		if len(msgs) > 0 {
			if prev, next := msgs[len(msgs)-1], msg; !prev.Time.Before(next.Time) {
				t.Errorf("invalid resut: prev %+v, next %+v", prev, next)
			}
		}
		msgs = append(msgs, msg)
	}

	if len(msgs) < 3 {
		t.Errorf("proceeded messages too few: %+v", msgs)
	}
}
