package redisz

import (
	"github.com/go-redis/redis"
	"time"
)

type Session struct {
	client *redis.Client
}

func NewSession(c Conf) Session {
	return Session{
		client: newRedis(c),
	}
}

func (s *Session) SetWithExpireTime(key string, val string, expireT time.Duration) {
	s.client.Set(key, val, expireT)
}

func (s *Session) Set(key string, val string) {
	s.client.Set(key, val, -1)
}

func (s *Session) Get(key string) (string, error) {
	ret := s.client.Get(key)
	return ret.Result()
}
