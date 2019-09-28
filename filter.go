package redisync

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func NewScoreFilter(pool Pool, opts ...Option) *ScoreFilter {
	return &ScoreFilter{
		Config: createConfig(opts),
		pool:   pool,
	}
}

type ScoreFilter struct {
	Config
	pool Pool
}

func (f *ScoreFilter) Filter(ctx context.Context, key string, score int) (ok bool, err error) {
	conn, err := f.pool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	defer func() {
		if err != nil {
			conn.Do("DISCARD")
		}
	}()

	_, err = conn.Do("WATCH", key)
	if err != nil {
		return false, err
	}

	got, err := redis.Int(conn.Do("GET", key))
	if (err != nil && err != redis.ErrNil) || (err == nil && got >= score) {
		return false, err
	}

	err = conn.Send("MULTI")
	if err != nil {
		return false, err
	}
	err = conn.Send("SET", key, score)
	if err != nil {
		return false, err
	}
	resp, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	if len(resp) < 1 || resp[0] != "OK" {
		return false, nil
	}

	return true, nil
}
