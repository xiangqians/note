// app
package typ

// Arg 应用参数
type Arg struct {
	Port     int    // 监听端口
	DataDir  string // 数据目录
	AllowReg int    // 是否允许用户注册，0-不允许，1-允许
}
