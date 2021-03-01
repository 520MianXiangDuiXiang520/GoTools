package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sync"
	"testing"
	"time"
)

type DBConn struct {
	Engine    string        `json:"engine"`
	DBName    string        `json:"db_name"`
	User      string        `json:"user"`
	Password  string        `json:"password"`
	Host      string        `json:"host"`
	Port      int           `json:"port"`
	MIdleConn int           `json:"max_idle_conn"` // 最大空闲连接数
	MOpenConn int           `json:"max_open_conn"` // 最大打开连接数
	MLifetime time.Duration `json:"max_lifetime"`  // 连接超时时间
	LogMode   bool          `json:"log_mode"`
}

type Table struct {
	gorm.Model
	Name string
}

func TestInitDBSetting(t *testing.T) {
	conn := DBConn{
		Engine:    "mysql",
		DBName:    "test",
		User:      "test",
		Password:  "123456",
		Host:      "127.0.0.1",
		Port:      3306,
		MIdleConn: 5,
		MOpenConn: 10,
		MLifetime: time.Second * 3,
		LogMode:   false,
	}
	err := InitDBSetting(&conn)
	if err != nil {
		t.Error(err)
	}
	// db := GetDB()
	// db.CreateTable(&Table{})
}

func doOnce(idle, open, num int) {
	wg := sync.WaitGroup{}
	conn := DBConn{
		Engine:    "mysql",
		DBName:    "test",
		User:      "test",
		Password:  "123456",
		Host:      "127.0.0.1",
		Port:      3306,
		MIdleConn: idle,
		MOpenConn: open,
		MLifetime: time.Second * 3,
		LogMode:   false,
	}
	wg.Add(num)
	for i := 0; i < num; i++ {
		_ = InitDBSetting(&conn)
		go func(i int) {

			db := GetDB()
			if pingErr := db.DB().Ping(); pingErr != nil {
				fmt.Println(pingErr)
			}
			db.Create(&Table{
				Name: fmt.Sprintf("%d", i),
			})

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func getTime(idle, open, num int) {
	var res int64
	for i := 0; i < 6; i++ {
		begin := time.Now().UnixNano()
		doOnce(idle, open, num)
		res = res + time.Now().UnixNano() - begin
		fmt.Println(res)
	}
	fmt.Printf("idle: %d, open: %d, --- %v ms \n", idle, open, res/6/1000/1000)
}

// func TestInitDBSetting2(t *testing.T) {
//
// 	getTime(10, 20, 500)
// 	// getTime(20, 40, 500)
// 	// getTime(40, 80, 500)
// 	// getTime(80, 160, 500)
// 	// getTime(160, 320, 500)
// 	// getTime(320, 640, 500)
// 	// getTime(200, 320, 500)
// }
