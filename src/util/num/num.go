// num
// @author xiangqian
// @date 12:32 2023/02/04
package num

import (
	"math/rand"
	"reflect"
	"time"
)

// RandIntn 获取 [0, n) 随机数
func RandIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func Int64(i any) int64 {
	//if v, r := i.(int); r {
	//	return int64(v)
	//}

	switch reflect.ValueOf(i).Kind() {
	case reflect.Int:
		return int64(i.(int))

	case reflect.Int8:
		return int64(i.(int8))

	case reflect.Int16:
		return int64(i.(int16))

	case reflect.Int32:
		return int64(i.(int32))

	case reflect.Int64:
		return i.(int64)

	default:
		return -1
	}
}
