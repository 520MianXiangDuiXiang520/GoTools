package dao

import (
	"errors"
	"log"
	"reflect"
)

type DaoLogic interface{}

// 将 def 以一个事务的方式执行。
//
// def 是执行一组 ORM 语句的函数，他应该满足以下条件：
//
// 1. 第一个参数是 *gorm.DB 类型，def 中的所有 ORM 操作都应该使用改对象。
//
// 2. 返回应该至少有一个是 error 类型的，如果有多个返回值时，error 类型的应该作为最后一个。
//
// 函数会返回 def 执行的结果，他以 reflect.Value 切片的形式返回
func UseTransaction(def DaoLogic, args []interface{}, logger *log.Logger) ([]reflect.Value, error) {
	var err error
	tx := GetDB().Begin()
	tx.LogMode(true)
	defer func() {
		// def 抛出 panic, 回滚
		if pan := recover(); pan != nil {
			logger.Printf("Transaction execution failed and has been rolled back！error: %v", pan)
			tx.Rollback()
		}
		// def 返回了一个 err, 回滚
		if err != nil {
			logger.Printf("Transaction return false and has been rolled back！error: %v", err)
			tx.Rollback()
		}
		tx.Commit()
	}()
	value := reflect.ValueOf(def)
	if value.Kind() != reflect.Func {
		return nil, errors.New("TypeError: def is not a Func type")
	}
	if reflect.TypeOf(args[0]) != reflect.TypeOf(tx) {
		return nil, errors.New("TypeError: the first parameter must be of type *DB")
	}
	argsVal := make([]reflect.Value, len(args))
	for i, arg := range args {
		argsVal[i] = reflect.ValueOf(arg)
	}

	argsVal[0] = reflect.ValueOf(tx)
	res := value.Call(argsVal)

	// 如果 def 主动抛出异常，回滚
	errVal := res[len(res)-1]

	if errVal.Interface() != nil {
		if _, ok := errVal.Interface().(error); ok {
			err = errVal.Interface().(error)
		} else {
			err = errors.New("return error")
		}
		return res, err
	}
	return res, nil
}
