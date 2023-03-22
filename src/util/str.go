// str
// @author xiangqian
// @date 12:31 2023/02/04
package util

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"strconv"
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

func TypeAsStr(i any) string {
	if i == nil {
		return ""
	}

	if v, r := i.(error); r {
		return v.Error()
	}

	return fmt.Sprintf("%v", i)
}

func StrAsType[T any](value string) (T, error) {
	var t T
	rflVal := reflect.ValueOf(t)
	//log.Println(rflVal)
	switch rflVal.Type().Kind() {
	case reflect.Int64:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(id).(T), err

	case reflect.String:
		return any(value).(T), nil
	}

	return t, errors.New(fmt.Sprintf("Unsupported operation: %v", rflVal.Type().Kind()))
}
