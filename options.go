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
	DefaultOnceExpiration = 3 * 24 * time.Second
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

func createDefaultOnceConfig() OnceConfig {
	return OnceConfig{
		Expiration: DefaultOnceExpiration,
	}
}

type OnceConfig struct {
	Expiration       time.Duration
	UnlockAfterError bool
}

type OnceOption func(*OnceConfig)

func createOnceConfig(opts []OnceOption) OnceConfig {
	c := createDefaultOnceConfig()
	for _, f := range opts {
		f(&c)
	}
	return c
}

func WithOnceExpiration(d time.Duration) OnceOption {
	return func(c *OnceConfig) { c.Expiration = d }
}

func WithOnceUnlockAfterError() OnceOption {
	return func(c *OnceConfig) { c.UnlockAfterError = true }
}
