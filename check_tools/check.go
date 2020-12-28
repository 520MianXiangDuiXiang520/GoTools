// 通过标签自动检查请求格式是否正确
package check_tools

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	checkLengthRegexp, _ = regexp.Compile("len:\\[\\d*,\\d*\\]")
	checkSizeRegexp, _   = regexp.Compile("size:\\[\\d*,\\d*\\]")
	checkMoreRegexp, _   = regexp.Compile("more:\\d*")
	checkLessRegexp, _   = regexp.Compile("less:\\d*")
	checkEqualRegexp, _  = regexp.Compile("less:\\d*")
)

func findNum(str string) []int {
	com := regexp.MustCompile("[0-9\\-]+")
	resS := com.FindAllString(str, -1)
	res := make([]int, len(resS))
	for i, r := range resS {
		vInt, err := strconv.Atoi(r)
		if err != nil {
			panic("WrongLabel:" + r)
		}
		res[i] = vInt
	}
	return res
}

func checkLen(tag, v string) bool {
	n := findNum(tag)
	if len(n) < 2 {
		panic("WrongLabel" + tag)
	}
	if !checkStringLen(v, n[0], n[1]) {
		log.Printf("[Check] %s length is %d, not in [%d, %d]\n", v, len(v), n[0], n[1])
		return false
	}
	return true
}

// 检查 int 类型的字段，包括(int, int8, int16, int32, int64)，大小写，空格不敏感
// - 检查不为 0: `check:"not null"` 或 `check:"not zero"`
// - 范围检查:   `check:"size: [0, 10]"`
func checkInt(v reflect.Value, tag string) bool {
	value := v.Int()
	tags := strings.Split(tag, ";")
	for _, t := range tags {
		m := strings.Replace(t, " ", "", -1)
		m = strings.ToLower(m)
		switch {
		case m == "notnull" || m == "notzero":
			if value == 0 {
				log.Printf("[Check] 0 is not null\n")
				return false
			}
		case checkSizeRegexp.Match([]byte(m)):
			n := findNum(m)
			if len(n) < 2 {
				panic("WrongLabel" + t)
			}
			if value < int64(n[0]) || value > int64(n[1]) {
				log.Printf("[Check]  %d not in [%d, %d]\n", value, n[0], n[1])
				return false
			}
		case checkMoreRegexp.Match([]byte(m)):
			// more:10
			num, err := strconv.Atoi(m[5:])
			if err != nil {
				panic("WrongLabel" + t)
			}
			if value <= int64(num) {
				log.Printf("[Check]  %d no greater than %d \n", value, num)
				return false
			}
		case checkLessRegexp.Match([]byte(m)):
			// less:10
			num, err := strconv.Atoi(m[5:])
			if err != nil {
				panic("WrongLabel" + t)
			}
			if value >= int64(num) {
				log.Printf("[Check]  %d not less than %d \n", value, num)
				return false
			}
		case checkEqualRegexp.Match([]byte(m)):
			// equal:10
			num, err := strconv.Atoi(m[6:])
			if err != nil {
				panic("WrongLabel" + t)
			}
			if value != int64(num) {
				log.Printf("[Check]  %d not equal to %d \n", value, num)
				return false
			}
		}
	}
	return true
}

// 检查 slice 类型
// 长度检查: `check:"len: [1, 10]"`
func checkSlice(v reflect.Value, tag string) bool {
	tags := strings.Split(tag, ";")
	for _, t := range tags {
		m := strings.Replace(t, " ", "", -1)
		m = strings.ToLower(m)
		switch {
		case checkLengthRegexp.Match([]byte(m)):
			n := findNum(tag)
			if len(n) < 2 {
				panic("WrongLabel" + tag)
			}
			if v.Len() < n[0] || v.Len() > n[1] {
				log.Printf("[Check] %s length is %d, not in [%d, %d]\n", v, v.Len(), n[0], n[1])
				return false
			}
		}
	}
	return true
}

func checkPtr(v reflect.Value, tag string) bool {
	tags := strings.Split(tag, ";")
	for _, t := range tags {
		m := strings.Replace(t, " ", "", -1)
		m = strings.ToLower(m)
		switch {
		case m == "notnull" || m == "notnil":
			if v.IsNil() {
				log.Printf("[Check] this pointer is nil")
				return false
			}
		}
	}
	return true
}

func switchType(field reflect.StructField, value reflect.Value, tag string) bool {
	switch field.Type.Kind() {
	case reflect.String:
		if !checkString(value, tag) {
			log.Printf("[check] [%v] Failed label inspection", field.Name)
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !checkInt(value, tag) {
			log.Printf("[check] [%v] Failed label inspection", field.Name)
			return false
		}
	case reflect.Slice:
		if !checkSlice(value, tag) {
			log.Printf("[check] [%v] Failed label inspection", field.Name)
			return false
		}
	case reflect.Ptr:
		if !checkPtr(value, tag) {
			log.Printf("[check] [%v] Failed label inspection", field.Name)
			return false
		}
	case reflect.Struct:
		if !checkStructReq(value.Interface()) {
			log.Printf("[check] [%v] Failed label inspection", field.Name)
			return false
		}
	}
	return true
}

func checkStructReq(s interface{}) bool {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("check")
		if !switchType(field, value, tag) {
			return false
		}
	}
	return true
}

func checkPtrReq(s interface{}) bool {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	if t.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("check: %v can not used in CheckRequest", t.Elem().Kind()))
	}
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := v.Elem().Field(i)
		tag := field.Tag.Get("check")
		if !switchType(field, value, tag) {
			return false
		}
	}
	return true
}

// 请使用 Check 方法
func CheckRequest(req interface{}) (ok bool) {
	t := reflect.TypeOf(req)
	switch t.Kind() {
	case reflect.Ptr:
		return checkPtrReq(req)
	case reflect.Struct:
		return checkStructReq(req)
	}
	panic(fmt.Sprintf("check: %v can not used in CheckRequest", t.Kind()))
	return false
}

func Check(req interface{}) (ok bool) {
	return CheckRequest(req)
}
