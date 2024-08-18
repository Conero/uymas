package str

import (
	"strings"
)

// @Date：   2018/11/7 0007 11:38
// @Author:  Joshua Conero
// @Name:    字符互队列

// InQuei checkout substring exist in array that case insensitive
func InQuei(s string, que []string) int {
	idx := -1
	s = strings.ToLower(s)
	for i, v := range que {
		if s == strings.ToLower(v) {
			idx = i
			break
		}
	}
	return idx
}

// StrQueueToAny string slice convert to any slice
func StrQueueToAny(args []string) []any {
	var anyQueue []any
	for _, s := range args {
		anyQueue = append(anyQueue, s)
	}
	return anyQueue
}
