package strings

import (
    `strings`
)

// DeDuplication 用于字符串去重
// 基于 map, 不能在并发场景下使用
func DeDuplication(s string) string {
    buf := strings.Builder{}
    dict := make(map[rune]struct{})
    for _, c := range s {
        if _, ok := dict[c]; !ok {
            buf.WriteRune(c)
            dict[c] = struct{}{}
        }
    }
    return buf.String()
}
