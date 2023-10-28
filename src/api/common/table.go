// table
// @author xiangqian
// @date 15:40 2023/10/28
package common

import (
	"reflect"
	"strings"
)

func Table[T any]() string {
	// 获取泛型类型
	var t T
	tType := reflect.TypeOf(t)
	// 结构体名称（此处即数据表名）
	return strings.ToLower(tType.Name())
}
