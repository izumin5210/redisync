package redisync

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func NewScoreFilter(pool Pool, m *Monitor, opts ...Option) *ScoreFilter {
	return &ScoreFilter{
		Config: createConfig(opts),
		m:      m,
		pool:   pool,
	}
}

type ScoreFilter struct {
	Config
	m    *Monitor
	pool Pool
}

func (f *ScoreFilter) Filter(ctx context.Context, key string, score int) (ok bool, err error) {
	err = f.m.Synchronize(ctx, key+":lock", func(ctx context.Context) (err error) {
		ok, err = f.filter(ctx, key, score)
		return
	})

	return
}

func (f *ScoreFilter) filter(ctx context.Context, key string, score int) (ok bool, err error) {
	conn, err := f.pool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	got, err := redis.Int(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return false, err
	}
	if err == redis.ErrNil || got < score {
		ok = true
		_, err = conn.Do("SET", key, score)
		if err != nil {
			return false, err
		}
	}

	return ok, nil
}
