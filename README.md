# GinTools

## check_tools

一个通过标签快速检查传入参数的工具：

使用方法：

```go
package main

import (
    `fmt`
    `github.com/520MianXiangDuiXiang520/GinTools/check_tools`
)

type DemoFather struct {
    FName string `check:"not null"`
}

type Demo struct {
    DemoFather
    Name string `check:"not null; len:[0, 12];"`
    Age  int    `check:"not null; size: [1, 150]"`
    Mail string `check:"not null; email"`
}

func main() {
    req := &Demo{}
    req.FName = "12"
    req.Name  = "1"
    req.Age   = 10
    req.Mail  = "15364968962@163,com"
    if check_tools.CheckRequest(req) {
        fmt.Printf("pass")
    }
}

// 2020/11/21 15:07:39 [Check] 15364968962@163,com is not a email
// 2020/11/21 15:07:39 [check] [Mail] Failed label inspection
```

## dao_tools

* 数据库连接工具（MySQL）
* 数据库事务工具

使用方法：

```go
func init() {
    daoUtils.InitDBSetting(src.GetSetting().Database,10, 30, time.Second*100, true)
}

func InsertToken(user *User, token string) (ok bool) {
    // 使用事务，保证一致性
    _, err := daoUtils.UseTransaction(func(db *gorm.DB, user *User, token string) (err error) {
        err = deleteTokenByUser(db, user)
        if err != nil {
            return
        }
        return insertToken(db, user, token)
    }, []interface{}{&gorm.DB{}, user, token })
    
	if err != nil {
        msg := fmt.Sprintf("Fail to insert token; user = %v, token = %v", user, token)
        utils.ExceptionLog(err, msg)
        return false
    }
	return true
}
```