package dao_tools

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)

var redisPool *redis.Pool
var redisLock sync.Mutex

type RedisConn struct {
	Host      string        `json:"host" check:"not null"`
	Password  string        `json:"password" check:"not null"`
	Port      int           `json:"port" check:"not null"`
	MIdleConn int           `json:"max_idle_conn"` // 最大空闲连接数
	MOpenConn int           `json:"max_open_conn"` // 最大打开连接数
	MLifetime time.Duration `json:"max_lifetime"`  // 连接超时时间
}

func (conn *RedisConn) newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conn.MIdleConn,
		MaxActive:   conn.MOpenConn,
		IdleTimeout: conn.MLifetime,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				conn.Host+":"+strconv.Itoa(conn.Port),
				redis.DialPassword(conn.Password),
			)
		},
	}
}

func initRedisSetting(from interface{}) (*RedisConn, error) {
	set := RedisConn{}
	err := transform(from, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}

// InitRedisPool 用来初始化一个 Redis 连接池 对象
// 你需要传递一个包含上面 RedisConn 结构体中字段的对象
// （无所谓类型，只需要 json 标签与类型对应即可，其中
// check 标签标记为 not null 的是必须字段）在应用程序中
// 连接池对象只需要初始化一遍即可使用 GetRedisConn 获取
// 连接实例去执行了，多次调用也不会多次实例化
func InitRedisPool(setting interface{}) error {
	if redisPool == nil {
		redisLock.Lock()
		if redisPool == nil {
			res, err := initRedisSetting(setting)
			if err != nil {
				return err
			}
			redisPool = res.newPool()
		}
		redisLock.Unlock()
	}
	return nil
}

// GetRedisConn 用来获取一个 redis 连接实例
// 使用之前应该确保连接池已经被初始化了
func GetRedisConn() redis.Conn {
	if redisPool == nil {
		panic("redis configuration not loaded！")
	}
	return redisPool.Get()
}
