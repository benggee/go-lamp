package mysqlz

import (
	"encoding/json"
	"fmt"
	"testing"
)


type User struct {
	Id int64 `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Realname string `db:"realname" json:"realname"`
	Age int `db:"age" json:"age"`
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
}

func TestA(t *testing.T) {
	s := NewSession("test:test@tcp(127.0.0.1:3306)/test?tls=skip-verify&autocommit=true", 1000 * 1000)
	u := make([]*User, 0)
	s.Rows(&u, "select * from user")


	ret1, err := json.Marshal(u)
	fmt.Println(string(ret1), err)
}