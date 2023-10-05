// 应用参数
// @author xiangqian
// @date 22:45 2023/06/12
package typ

import (
	"fmt"
)

var arg Arg

// Arg 应用参数
type Arg struct {
	TimeZone string // 时区，如：Asia/Shanghai（上海时区）
	Port     int    // 监听端口
	Path     string // 服务根路径
	DataDir  string // 数据目录
	AllowReg bool   // 是否允许用户注册
}

// String 返回结构体类型字符串
func (arg Arg) String() string {
	return fmt.Sprintf("Arg { TimeZone = %s, Port = %d, Path = %s, DataDir = %s, AllowReg = %v }", arg.TimeZone, arg.Port, arg.Path, arg.DataDir, arg.AllowReg)
}

func GetArg() Arg {
	return arg
}

func SetArg(_arg Arg) {
	arg = _arg
}
