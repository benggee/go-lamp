package redisz

type Conf struct {
	Addr     []string `json:"addr"`
	Password string   `json:"password"`
}
