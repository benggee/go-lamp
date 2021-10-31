package redisz

import "github.com/go-redis/redis"

type Session struct {
	client *redis.Client
}

func NewSession(c Conf) Session {
	return Session{
		client: newRedis(c),
	}
}
