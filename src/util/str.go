// str
// @author xiangqian
// @date 12:31 2023/02/04
package util

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
)

// Uuid https://github.com/google/uuid
func Uuid() string {
	return uuid.New().String()
}

// NameHumpToUnderline 驼峰转下划线
func NameHumpToUnderline(name string) string {
	pRegexp := regexp.MustCompile("([A-Z])")
	r := pRegexp.FindAllIndex([]byte(name), -1)
	l := len(r)
	if l == 0 {
		return strings.ToLower(name)
	}

	var res = make([]string, l+1)
	resIdx := 0
	index := 0
	for _, v := range r {
		s := name[index:v[0]]
		if s != "" {
			res[resIdx] = s
			resIdx++
		}
		index = v[0]
	}
	res[resIdx] = name[index:]
	for i, s := range res {
		if s == "" {
			res = res[0:i]
			break
		}
	}
	return strings.ToLower(strings.Join(res, "_"))
}
