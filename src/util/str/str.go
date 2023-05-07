// str
// @author xiangqian
// @date 12:31 2023/02/04
package str

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

// ConvNameHumpToUnderline 驼峰转下划线
func ConvNameHumpToUnderline(name string) string {
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

// ConvTypeToStr 类型转string
func ConvTypeToStr(i any) string {
	if i == nil {
		return ""
	}

	if v, r := i.(error); r {
		return v.Error()
	}

	return fmt.Sprintf("%v", i)
}

// ConvStrToType string转类型（基本数据类型）
func ConvStrToType[T any](value string) (T, error) {
	var t T
	rflVal := reflect.ValueOf(t)
	//log.Println(rflVal)
	switch rflVal.Type().Kind() {
	case reflect.Int:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(int(any(id).(int64))).(T), err

	case reflect.Int8:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(int8(any(id).(int64))).(T), err

	case reflect.Uint8:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(uint8(any(id).(int64))).(T), err

	case reflect.Int64:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(id).(T), err

	case reflect.String:
		return any(value).(T), nil
	}

	return t, errors.New(fmt.Sprintf("Unsupported operation: %v", rflVal.Type().Kind()))
}
