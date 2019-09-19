package redisync_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/izumin5210/redisync"
)

func TestMutex(t *testing.T) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.DialURL(os.Getenv("REDIS_URL")) },
	}
	defer pool.Close()

	defer func() {
		conn := pool.Get()
		defer conn.Close()
		conn.Do("FLUSHALL")
	}()

	ctx := context.Background()
	m1 := redisync.NewMutex(pool, "mutex1")
	m2 := redisync.NewMutex(pool, "mutex2")

	var (
		cnt int
		wg  sync.WaitGroup
	)

	err := m1.Lock(ctx)
	if err != nil {
		t.Errorf("Lock() returned an error: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := m1.Lock(ctx)
		if err != nil {
			t.Errorf("Lock() returned an error: %v", err)
		}
		if got, want := cnt, 1; got != want {
			t.Errorf("count is %d, want %d", got, want)
		}
		cnt++
	}()

	time.Sleep(50 * time.Millisecond)
	cnt++

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(20 * time.Millisecond)
		err := m2.Lock(ctx)
		if err != nil {
			t.Errorf("Lock() returned an error: %v", err)
		}
		if got, want := cnt, 3; got != want {
			t.Errorf("count is %d, want %d", got, want)
		}
		cnt++

		err = m2.Unlock(ctx)
		if err != nil {
			t.Errorf("Unlock() returned an error: %v", err)
		}
	}()

	err = m2.Lock(ctx)
	if err != nil {
		t.Errorf("Lock() returned an error: %v", err)
	}
	err = m1.Unlock(ctx)
	if err != nil {
		t.Errorf("Unlock() returned an error: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	if got, want := cnt, 2; got != want {
		t.Errorf("count is %d, want %d", got, want)
	}
	cnt++

	err = m1.Unlock(ctx)
	if err != nil {
		t.Errorf("Unlock() returned an error: %v", err)
	}

	err = m2.Unlock(ctx)
	if err != nil {
		t.Errorf("Unlock() returned an error: %v", err)
	}

	wg.Wait()
}

func TestMutex_WithOptions(t *testing.T) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.DialURL(os.Getenv("REDIS_URL")) },
	}
	defer pool.Close()

	defer func() {
		conn := pool.Get()
		defer conn.Close()
		conn.Do("FLUSHALL")
	}()

	ctx := context.Background()
	m := redisync.NewMutex(pool, "mutex",
		redisync.WithLockExpiration(1*time.Second),
		redisync.WithBackOffFactory(redisync.BackOffFactoryFunc(func() backoff.BackOff {
			return backoff.WithMaxRetries(backoff.NewConstantBackOff(10*time.Millisecond), 3)
		})),
	)

	err := m.Lock(ctx)
	if got, want := err, error(nil); got != want {
		t.Errorf("Lock returned %v, want %v", got, want)
	}

	err = m.Lock(ctx)
	if got, want := err, redisync.ErrConflict; got != want {
		t.Errorf("Lock returned %v, want %v", got, want)
	}

	time.Sleep(1 * time.Second)

	err = m.Lock(ctx)
	if got, want := err, error(nil); got != want {
		t.Errorf("Lock returned %v, want %v", got, want)
	}
}
