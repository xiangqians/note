// 应用参数
// @author xiangqian
// @date 22:45 2023/06/12
package model

import (
	"fmt"
)

var arg Arg

// Arg 应用参数
type Arg struct {
	TimeZone    string // 时区，如：Asia/Shanghai（上海时区）
	Port        int    // 监听端口
	ContextPath string // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	DataDir     string // 数据目录
	AllowSignUp bool   // 是否允许用户注册
}

// String 返回结构体类型字符串
func (arg Arg) String() string {
	return fmt.Sprintf("Arg { TimeZone = %s, Port = %d, ContextPath = %s, DataDir = %s, AllowSignUp = %v }", arg.TimeZone, arg.Port, arg.ContextPath, arg.DataDir, arg.AllowSignUp)
}

func GetArg() Arg {
	return arg
}

func SetArg(_arg Arg) {
	arg = _arg
}
