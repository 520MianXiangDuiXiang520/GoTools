package check_tools

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

func checkStringLen(v string, min, max int) bool {
	return utf8.RuneCountInString(v) >= min && utf8.RuneCountInString(v) <= max
}

func isEmail(v string) bool {
	if len(v) <= 0 {
		return false
	}
	emailRules := `^[A-Za-z0-9_.]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	ok, err := regexp.MatchString(emailRules, v)
	return ok && err == nil
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
				log.Printf("[Check] \"\" is not null\n")
				return false
			}
		case m == "email":
			if !isEmail(v) {
				log.Printf("[Check] %s is not a email\n", v)
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
