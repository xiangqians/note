// app type
// @author xiangqian
// @date 20:53 2023/03/13
package typ

import (
	"fmt"
)

// AppArg 应用参数
type AppArg struct {
	Loc      string // location，如：Asia/Shanghai（上海时区）
	Port     int    // 监听端口
	DataDir  string // 数据目录
	AllowReg int    // 是否允许用户注册，0-不允许，1-允许
}

// String 返回结构体类型字符串
func (appArg AppArg) String() string {
	return fmt.Sprintf("AppArg { Loc = %s, Port = %d, DataDir = %s, AllowReg = %d }", appArg.Loc, appArg.Port, appArg.DataDir, appArg.AllowReg)
}
