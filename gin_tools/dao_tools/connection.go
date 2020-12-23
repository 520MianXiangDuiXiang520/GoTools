package dao_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/check_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
	"sync"
	"time"
)

// 定义建立数据库连接时所需要的值，可以使用任意类型的 struct,
// 只要 json 标签于此对于即可，对于类似 Engine, User 等必须参数
// 我们会使用 check_tools.Check 检查，请确保值正确。
type DBConnector struct {
	Engine      string        `json:"engine" check:"not null"`
	DBName      string        `json:"db_name" check:"not null"`
	User        string        `json:"user" check:"not null"`
	Password    string        `json:"password" check:"not null"`
	Host        string        `json:"host" check:"not null"`
	Port        int           `json:"port" check:"not null"`
	MaxIdleConn int           `json:"max_idle_conn"` // 最大空闲连接数
	MaxOpenConn int           `json:"max_open_conn"` // 最大打开连接数
	MaxLifetime time.Duration `json:"max_lifetime"`  // 连接超时时间
	LogMode     bool          `json:"log_mode"`
}

func (conn *DBConnector) NewConnect() *gorm.DB {
	connURI := ""
	switch strings.ToLower(conn.Engine) {
	case "mysql":
		connURI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			conn.User, conn.Password,
			conn.Host, conn.Port, conn.DBName)
	case "":
		panic("engine is nil")
	default:
		panic("unrecognized database engine")
	}

	db, err := gorm.Open(conn.Engine, connURI)
	if err != nil {
		msg := fmt.Sprintf("Fail to open db, URI is %s", connURI)
		utils.ExceptionLog(err, msg)
		panic(err)
	}
	return db
}

// 设置连接池参数
func setup(maxIdle, maxOpen int, maxLifeTime time.Duration, logMode bool) {
	db.DB().SetMaxOpenConns(maxOpen)        // 最大连接数
	db.DB().SetMaxIdleConns(maxIdle)        // 最大空闲连接数
	db.DB().SetConnMaxLifetime(maxLifeTime) // 设置连接空闲超时
	db.LogMode(logMode)
}

var (
	dbConnector *DBConnector
	db          *gorm.DB
	lock        = sync.Mutex{}
)

func initDBSetting(set interface{}) (*DBConnector, error) {
	conn := DBConnector{}
	data, err := json.Marshal(set)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("marshal %v Fail!", set))
		return nil, err
	}
	err = json.Unmarshal(data, &conn)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to UnMarchall"))
		return nil, err
	}
	if !check_tools.Check(conn) {
		return nil, errors.New("miss field")
	}
	return &conn, nil
}

// 初始化数据库连接.
//
// 对于 set, 可以传入一个任意类型的 struct, 但是需要保证 set 里包含
// 初始化连接所必须的参数，并且使用正确的 json 标签，具体参数和标签可以
// 参考 DBConnector, 其中 check 标签标记为 not null 的是必须参数,如:
//    conn := DBConn{
//        Engine: "mysql",
//        DBName: "test",
//        User: "test",
//        Password: "123456",
//        Host: "127.0.0.1",
//        Port: 3306,
//        MIdleConn: idle,
//        MOpenConn: open,
//        MLifetime: time.Second * 3,
//        LogMode: false,
//    }
//    err := InitDBSetting(&conn)
//
// 若初始化成功，则返回值 error 应该为 nil。
//
// 初始化成功后，可以调用 GetDB 获取具体的 *gorm.DB 对象.
func InitDBSetting(set interface{}) error {
	if dbConnector == nil {
		lock.Lock()
		if dbConnector == nil {
			res, err := initDBSetting(set)
			if err != nil {
				return err
			}
			dbConnector = res
		}
		lock.Unlock()
	}
	return nil
}

// 获取一个 *gorm.DB 对象，使用前必须使用 InitDBSetting 初始化连接
// 否则会导致 panic
func GetDB() *gorm.DB {
	if !check_tools.Check(dbConnector) {
		panic("Database configuration is not loaded")
	}
	if db == nil {
		lock.Lock()
		if db == nil {
			db = dbConnector.NewConnect()
			setup(dbConnector.MaxIdleConn, dbConnector.MaxOpenConn, dbConnector.MaxLifetime, dbConnector.LogMode)
		}
		lock.Unlock()
	}
	return db
}
