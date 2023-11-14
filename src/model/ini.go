// @author xiangqian
// @date 22:00 2023/11/02
package model

import (
	"fmt"
	util_os "note/src/util/os"

	// https://pkg.go.dev/gopkg.in/ini.v1
	// https://github.com/go-ini/ini
	pkg_ini "gopkg.in/ini.v1"

	"time"
)

// 配置
type ini struct {
	Log    log    // 日志配置
	Db     db     // 数据库配置
	Data   data   // 数据配置
	Server server // 服务配置
}

// 日志配置
type log struct {
	Dir        string // 日志文件目录
	FileName   string // 日志文件名
	MaxSize    int64  // 日志文件大小，单位：字节
	MaxHistory int    // 日志文件最大历史记录
}

// 数据库配置
type db struct {
	Driver          string        // 驱动名
	Dns             string        // 数据源
	MaxOpenConns    int           // 池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限
	ConnMaxLifetime time.Duration // 一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用
	MaxIdleConns    int           // 池中“空闲”连接数的上限
	ConnMaxIdleTime time.Duration // 在被标记为失效之前一个连接最长空闲时间
}

// 数据配置
type data struct {
	Dir string // 物理数据文件目录
}

// 服务配置
type server struct {
	Port             uint16        // 监听端口
	ContextPath      string        // 应用的上下文路径，也可以称为项目路径，是构成url地址的一部分
	SessionSecretKey string        // 会话密钥
	SessionMaxAge    time.Duration // 会话过期时间
	OpenSignup       bool          // 是否开放注册功能
}

var Ini ini

func init() {
	source := "E:\\workspace\\goland\\note\\res\\note.ini"
	//source := "./note.ini"
	file, err := pkg_ini.Load(source)
	if err != nil {
		panic(err)
	}

	// log
	log, err := file.GetSection("log")
	if err != nil {
		panic(err)
	}
	Ini.Log.Dir = log.Key("dir").MustString("./log")
	Ini.Log.FileName = log.Key("file-name").MustString("debug.log")
	maxSize, err := util_os.ParseByte(log.Key("max-size").MustString("10MB"))
	if err != nil {
		panic(err)
	}
	Ini.Log.MaxSize = int64(maxSize)
	Ini.Log.MaxHistory = log.Key("max-history").MustInt(2)

	// db
	db, err := file.GetSection("db")
	if err != nil {
		panic(err)
	}
	Ini.Db.Driver = db.Key("driver").String()
	Ini.Db.Dns = db.Key("dns").String()
	Ini.Db.MaxOpenConns = db.Key("max-open-conns").MustInt(3)
	Ini.Db.ConnMaxLifetime = db.Key("conn-max-lifetime").MustDuration(0 * time.Minute)
	Ini.Db.MaxIdleConns = db.Key("max-idle-conns").MustInt(2)
	Ini.Db.ConnMaxIdleTime = db.Key("conn-max-idle-time").MustDuration(30 * time.Minute)

	// data
	data, err := file.GetSection("data")
	if err != nil {
		panic(err)
	}
	Ini.Data.Dir = data.Key("dir").String()

	// server
	server, err := file.GetSection("server")
	if err != nil {
		panic(err)
	}
	Ini.Server.Port = uint16(server.Key("port").MustUint64(8080))
	Ini.Server.ContextPath = server.Key("context-path").MustString("/")
	Ini.Server.SessionSecretKey = server.Key("session-secret-key").String()
	Ini.Server.SessionMaxAge = server.Key("session-max-age").MustDuration(12 * time.Hour)
	Ini.Server.OpenSignup = server.Key("open-signup").MustBool(true)
}

// String 返回结构体类型字符串
func (ini ini) String() string {
	logString := fmt.Sprintf("Log\t\t{ Dir = %s, MaxSize = %f, MaxHistory = %d }", Ini.Log.Dir, Ini.Log.MaxSize, Ini.Log.MaxHistory)
	dbString := fmt.Sprintf("Db\t\t{ Driver = %s, Dns = %s, MaxOpenConns = %d, ConnMaxLifetime = %s, MaxIdleConns = %d, ConnMaxIdleTime = %s }",
		ini.Db.Driver, ini.Db.Dns, ini.Db.MaxOpenConns, ini.Db.ConnMaxLifetime, ini.Db.MaxIdleConns, ini.Db.ConnMaxIdleTime)
	dataString := fmt.Sprintf("Data\t{ Dir = %s }", ini.Data.Dir)
	serverString := fmt.Sprintf("Server\t{ Port = %d, ContextPath = %s, SessionSecretKey = %s, SessionMaxAge = %s, OpenSignup = %t }", ini.Server.Port, ini.Server.ContextPath, ini.Server.SessionSecretKey, ini.Server.SessionMaxAge, ini.Server.OpenSignup)
	return fmt.Sprintf("Ini"+
		"\n\t%s"+
		"\n\t%s"+
		"\n\t%s"+
		"\n\t%s",
		logString,
		dbString,
		dataString,
		serverString)
}
