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
