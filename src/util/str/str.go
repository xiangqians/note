// str
// @author xiangqian
// @date 12:31 2023/02/04
package str

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

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

	return t, errors.New(fmt.Sprintf("This type does not support conversion: %v", rflVal.Type().Kind()))
}
