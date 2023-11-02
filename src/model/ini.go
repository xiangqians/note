// @author xiangqian
// @date 22:00 2023/11/02
package model

// https://pkg.go.dev/gopkg.in/ini.v1
// https://github.com/go-ini/ini
import (
	"fmt"
	pkg_ini "gopkg.in/ini.v1"
)

type ini struct {
	Sys    sys    // 系统配置
	Server server // 服务配置
	Db     db     // 系统配置
}

type sys struct {
	TimeZone        string // 时区，如：Asia/Shanghai（上海时区）
	Mode            string // 模式，debug/release
	Port            int    // 监听端口
	ContextPath     string // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	OpenSignup      bool   // 是否开放注册功能
	Driver          string // driver name
	Dns             string // data source name
	MaxOpenConns    int    // 设置池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限。
	ConnMaxLifetime int    // 设置一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用。
	MaxIdleConns    int    // 设置池中“空闲”连接数的上限。缺省情况下，最大空闲连接数为2。
	ConnMaxIdleTime int    // 在被标记为失效之前一个连接最长空闲时间。
	DataDir         string // 物理数据文件目录
}

type server struct {
	TimeZone        string // 时区，如：Asia/Shanghai（上海时区）
	Mode            string // 模式，debug/release
	Port            int    // 监听端口
	ContextPath     string // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	OpenSignup      bool   // 是否开放注册功能
	Driver          string // driver name
	Dns             string // data source name
	MaxOpenConns    int    // 设置池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限。
	ConnMaxLifetime int    // 设置一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用。
	MaxIdleConns    int    // 设置池中“空闲”连接数的上限。缺省情况下，最大空闲连接数为2。
	ConnMaxIdleTime int    // 在被标记为失效之前一个连接最长空闲时间。
	DataDir         string // 物理数据文件目录
}

type db struct {
	TimeZone        string // 时区，如：Asia/Shanghai（上海时区）
	Mode            string // 模式，debug/release
	Port            int    // 监听端口
	ContextPath     string // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	OpenSignup      bool   // 是否开放注册功能
	Driver          string // driver name
	Dns             string // data source name
	MaxOpenConns    int    // 设置池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限。
	ConnMaxLifetime int    // 设置一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用。
	MaxIdleConns    int    // 设置池中“空闲”连接数的上限。缺省情况下，最大空闲连接数为2。
	ConnMaxIdleTime int    // 在被标记为失效之前一个连接最长空闲时间。
	DataDir         string // 物理数据文件目录
}

type data struct {
	TimeZone        string // 时区，如：Asia/Shanghai（上海时区）
	Mode            string // 模式，debug/release
	Port            int    // 监听端口
	ContextPath     string // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	OpenSignup      bool   // 是否开放注册功能
	Driver          string // driver name
	Dns             string // data source name
	MaxOpenConns    int    // 设置池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限。
	ConnMaxLifetime int    // 设置一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用。
	MaxIdleConns    int    // 设置池中“空闲”连接数的上限。缺省情况下，最大空闲连接数为2。
	ConnMaxIdleTime int    // 在被标记为失效之前一个连接最长空闲时间。
	DataDir         string // 物理数据文件目录
}

var Ini ini

func init() {
	source := "E:\\workspace\\goland\\note\\res\\note.ini"
	//source := "./note.ini"
	file, err := pkg_ini.Load(source)
	if err != nil {
		panic(err)
	}

	// sys
	sys, err := file.GetSection("sys")
	if err != nil {
		panic(err)
	}
	Ini.TimeZone = sys.Key("time-zone").MustString("Asia/Shanghai")

	// server
	server, err := file.GetSection("server")
	if err != nil {
		panic(err)
	}
	Ini.TimeZone = sys.Key("time-zone").MustString("Asia/Shanghai")

}

// String 返回结构体类型字符串
func (ini ini) String() string {
	return fmt.Sprintf("Arg { TimeZone = %s, Port = %d, ContextPath = %s, DataDir = %s, AllowSignUp = %v }",
		arg.TimeZone, arg.Port, arg.ContextPath, arg.DataDir, arg.AllowSignUp)
}
