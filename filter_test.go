package redisync_test

import (
	"context"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/izumin5210/redisync"
)

func TestScoreFilter(t *testing.T) {
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
	filter := redisync.NewScoreFilter(pool, m)

	var (
		ok  bool
		err error
	)

	ok, err = filter.Filter(context.Background(), "foo", 100)
	if err != nil {
		t.Errorf("returned %v, want nil", err)
	}
	if got, want := ok, true; got != want {
		t.Errorf("returned %t, want %t", got, want)
	}

	ok, err = filter.Filter(context.Background(), "foo", 100)
	if err != nil {
		t.Errorf("returned %v, want nil", err)
	}
	if got, want := ok, false; got != want {
		t.Errorf("returned %t, want %t", got, want)
	}

	ok, err = filter.Filter(context.Background(), "foo", 120)
	if err != nil {
		t.Errorf("returned %v, want nil", err)
	}
	if got, want := ok, true; got != want {
		t.Errorf("returned %t, want %t", got, want)
	}

	ok, err = filter.Filter(context.Background(), "foo", 100)
	if err != nil {
		t.Errorf("returned %v, want nil", err)
	}
	if got, want := ok, false; got != want {
		t.Errorf("returned %t, want %t", got, want)
	}

	ok, err = filter.Filter(context.Background(), "foo", 200)
	if err != nil {
		t.Errorf("returned %v, want nil", err)
	}
	if got, want := ok, true; got != want {
		t.Errorf("returned %t, want %t", got, want)
	}
}
