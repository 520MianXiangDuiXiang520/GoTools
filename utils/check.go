// 通过标签自动检查请求格式是否正确
package utils

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

func checkStringLen(v string, min, max int) bool {
	return len(v) >= min && len(v) <= max
}

func isEmail(v string) bool {
	if len(v) <= 0 {
		return false
	}
	emailRules := `^[A-Za-z0-9_.]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	ok, err := regexp.MatchString(emailRules, v)
	return ok && err == nil
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

// 检查 string 类型的字段, 大小写，空格不敏感
// - 检查空字符串:  `check:"not null"`
// - 检查字符串长度: `check:"len: [0, 10]"`
// - 检查邮箱格式:   `check:"email"`
func checkString(value reflect.Value, tag string) (ok bool) {
	tags := strings.Split(tag, ";")
	v := value.String()
	for _, t := range tags {
		m := strings.Replace(t, " ", "", -1)
		m = strings.ToLower(m)
		switch {
		case m == "notnull":
			if v == "" {
				log.Flags()
				log.Printf("[Check] \"\" is not null\n")
				return false
			}
		case m == "email":
			if !isEmail(v) {
				log.Printf("[Check] %s is not a emai\n", v)
				return false
			}
		case checkLengthRegexp.Match([]byte(m)):
			if !checkLen(t, v) {
				return false
			}
		}

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
		case checkSizeRegexp.Match([]byte(t)):
			n := findNum(t)
			if len(n) < 2 {
				panic("WrongLabel" + t)
			}
			if value < int64(n[0]) || value > int64(n[1]) {
				log.Printf("[Check]  %d not in [%d, %d]\n", value, n[0], n[1])
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
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !checkInt(value, tag) {
			return false
		}
	case reflect.Slice:
		if !checkSlice(value, tag) {
			return false
		}
	case reflect.Ptr:
		if !checkPtr(value, tag) {
			return false
		}
	case reflect.Struct:
		if !checkStructReq(value.Interface()) {
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
