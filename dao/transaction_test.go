package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

func init() {
	conn := DBConn{
		Engine:    "mysql",
		DBName:    "t1",
		User:      "root",
		Password:  "1234567",
		Host:      "127.0.0.1",
		Port:      3306,
		MIdleConn: 5,
		MOpenConn: 10,
		MLifetime: time.Second * 3,
		LogMode:   false,
	}
	InitDBSetting(&conn)
}

func ExampleUseTransaction() {

	def := func(db *gorm.DB, res *Table, id uint, name string) (ok bool, err error) {
		err = db.Select("id = ?, name = ?", id, name).First(&res).Error
		if err != nil {
			return false, err
		}
		err = db.Create(&Table{
			Name: "p",
		}).Error
		if err != nil {
			return false, err
		}
		return true, nil
	}
	res := Table{}
	args := []interface{}{
		&gorm.DB{}, &res, uint(83), "13",
	}

	resL, err := UseTransaction(def, args, log.New(os.Stdout, "[ Transaction ] ", log.LstdFlags))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
	if resL[0].Bool() {
		fmt.Println("insert success")
	}
}

// func TestUseTransaction(t *testing.T) {
// 	ExampleUseTransaction()
// }
