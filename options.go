package redisync

import (
	"time"

	"github.com/cenkalti/backoff/v3"
)

var (
	DefaultBackOffFactory BackOffFactory = BackOffFactoryFunc(func() backoff.BackOff {
		bo := backoff.NewExponentialBackOff()
		bo.MaxInterval = 1 * time.Second
		bo.InitialInterval = 5 * time.Millisecond
		return bo
	})
	DefaultLockExpiration = 120 * time.Second
)

func createDefaultConfig() Config {
	return Config{
		BackOffFactory: DefaultBackOffFactory,
		LockExpiration: DefaultLockExpiration,
	}
}

func createConfig(opts []Option) Config {
	c := createDefaultConfig()
	for _, f := range opts {
		f(&c)
	}
	return c
}

type Config struct {
	BackOffFactory BackOffFactory
	LockExpiration time.Duration
}

type Option func(*Config)

func WithBackOffFactory(f BackOffFactory) Option {
	return func(c *Config) { c.BackOffFactory = f }
}

func WithLockExpiration(d time.Duration) Option {
	return func(c *Config) { c.LockExpiration = d }
}
