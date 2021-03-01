package dao

import (
	"log"
	"time"
)

var redisConn = struct {
	Host      string        `json:"host" check:"not null"`
	Password  string        `json:"password" check:"not null"`
	Port      int           `json:"port" check:"not null"`
	MIdleConn int           `json:"max_idle_conn"` // 最大空闲连接数
	MOpenConn int           `json:"max_open_conn"` // 最大打开连接数
	MLifetime time.Duration `json:"max_lifetime"`  // 连接超时时间
}{
	Host:      "127.0.0.1",
	Password:  "12345678",
	Port:      8100,
	MLifetime: 3000,
	MOpenConn: 10,
	MIdleConn: 15,
}

// func TestGetRedis(t *testing.T) {
// 	err := InitRedisPool(redisConn)
// 	if err != nil {
// 		t.Error("init pool fail")
// 	}
// 	redis := GetRedisConn()
// 	_, err = redis.Do("SET", "test", "test")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func ExampleGetRedis() {
	// 使用 GetRedisConn 前必须初始化连接池
	// 连接池只需要初始化一遍
	err := InitRedisPool(redisConn)
	if err != nil {
		log.Println("fail to init pool")
	}
	redis := GetRedisConn()
	_, err = redis.Do("SET", "test", "test")
	// ...
}
