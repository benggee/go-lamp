package mysqlz

import "time"

type Conf struct {
	Dns             string        `json:"dns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	MaxOpenConns    int           `json:"maxOpenConns"`
	MaxIdleConns    int           `json:"maxIdleConns"`
}
