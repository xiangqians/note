// num
// @author xiangqian
// @date 12:32 2023/02/04
package util

import (
	"math/rand"
	"time"
)

// RandIntn 获取 [0, n) 随机数
func RandIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func Int64(i any) int64 {
	if v, r := i.(int); r {
		return int64(v)
	}

	if v, r := i.(int8); r {
		return int64(v)
	}

	if v, r := i.(int16); r {
		return int64(v)
	}

	if v, r := i.(int32); r {
		return int64(v)
	}

	if v, r := i.(int64); r {
		return v
	}

	return 0
}
