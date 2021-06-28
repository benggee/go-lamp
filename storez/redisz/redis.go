package redisz

import (
	"github.com/go-redis/redis"
)

type Option func(opt *redis.Options)

func newRedis(c Conf, opt ...Option) *redis.Client {
	if len(c.Addr) == 0 {
		panic("redis conf invalid")
	}

	options := &redis.Options{
		Addr: c.Addr[0],
	}

	if len(c.Password) == 0 {
		options.Password = c.Password
	}

	return redis.NewClient(options)
}
