package redisync_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/izumin5210/redisync"
)

func TestOnce(t *testing.T) {
	defer cleanupTestRedis()

	ctx := context.Background()
	once1 := redisync.NewOnce(pool)
	once2 := redisync.NewOnce(pool)

	var (
		wg                   sync.WaitGroup
		fooCnt, barCnt       uint32
		fooErrCnt, barErrCnt uint32
	)
	for i := 0; i < 100; i++ {
		for _, tc := range []struct {
			once        *redisync.Once
			cnt, errCnt *uint32
			key         string
		}{
			{once: once1, cnt: &fooCnt, errCnt: &fooErrCnt, key: "foo"},
			{once: once2, cnt: &fooCnt, errCnt: &fooErrCnt, key: "foo"},
			{once: once1, cnt: &barCnt, errCnt: &barErrCnt, key: "bar"},
			{once: once2, cnt: &barCnt, errCnt: &barErrCnt, key: "bar"},
		} {
			tc := tc
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := tc.once.Do(ctx, tc.key, func(context.Context) error {
					atomic.AddUint32(tc.cnt, 1)
					return nil
				})
				if err != nil {
					if got, want := err, redisync.ErrLocked; got != want {
						t.Errorf("Run() returned unexpected error %v, want %v", got, want)
					} else {
						atomic.AddUint32(tc.errCnt, 1)
					}
				}
			}()
		}
	}

	wg.Wait()

	if got, want := fooCnt, uint32(1); got != want {
		t.Errorf("foo called %d times, want %d", got, want)
	}

	if got, want := barCnt, uint32(1); got != want {
		t.Errorf("bar called %d times, want %d", got, want)
	}

	if got, want := fooErrCnt, uint32(199); got != want {
		t.Errorf("foo skipped %d times, want %d", got, want)
	}

	if got, want := barErrCnt, uint32(199); got != want {
		t.Errorf("bar skipped %d times, want %d", got, want)
	}
}

func TestOnce_WithOnceUnlockAfterError(t *testing.T) {
	defer cleanupTestRedis()

	ctx := context.Background()
	once := redisync.NewOnce(pool, redisync.WithOnceUnlockAfterError())

	var callCnt int32

	_ = once.Do(ctx, "foo", func(context.Context) error {
		atomic.AddInt32(&callCnt, 1)
		return errors.New("errors")
	})

	_ = once.Do(ctx, "foo", func(context.Context) error {
		atomic.AddInt32(&callCnt, 1)
		return nil
	})

	_ = once.Do(ctx, "foo", func(context.Context) error {
		atomic.AddInt32(&callCnt, 1)
		return nil
	})

	if got, want := callCnt, int32(2); got != want {
		t.Errorf("Closure in redisync.Once is called %d times, want %d", got, want)
	}
}
