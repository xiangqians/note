// arg
// @author xiangqian
// @date 22:45 2023/06/12
package typ

import (
	"fmt"
)

// Arg 应用参数
type Arg struct {
	Loc      string // location，如：Asia/Shanghai（上海时区）
	Port     int    // 监听端口
	Path     string // 项目根路径
	DataDir  string // 数据目录
	AllowReg int    // 是否允许用户注册，0-不允许，1-允许
}

// String 返回结构体类型字符串
func (arg Arg) String() string {
	return fmt.Sprintf("AppArg { Loc = %s, Port = %d, Path = %s, DataDir = %s, AllowReg = %d }", arg.Loc, arg.Port, arg.Path, arg.DataDir, arg.AllowReg)
}
