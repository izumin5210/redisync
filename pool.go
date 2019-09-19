package redisync

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Pool interface {
	GetContext(context.Context) (redis.Conn, error)
}
