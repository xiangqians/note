// 应用参数
// @author xiangqian
// @date 19:47 2023/07/10
package app

import (
	"flag"
	"log"
	"note/src/typ"
	util_os "note/src/util/os"
	"os"
	"path/filepath"
	"strings"
)

var arg typ.Arg

// 初始化应用参数
func initArg() {
	var timeZone string
	var port int
	var path string
	var dataDir string
	var allowReg string

	// eg:
	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"

	// 解析参数
	flag.StringVar(&timeZone, "timeZone", "Asia/Shanghai", "-timeZone Asia/Shanghai")
	flag.IntVar(&port, "port", 8080, "-port 8080")
	flag.StringVar(&path, "path", "/", "-path /")
	flag.StringVar(&dataDir, "dataDir", "./data", "-dataDir ./data")
	flag.StringVar(&allowReg, "allowReg", "true", "-allowReg true")
	flag.Parse()

	// 时区
	timeZone = strings.TrimSpace(timeZone)

	// 服务根路径
	path = strings.TrimSpace(path)
	if path == "/" {
		path = ""
	}

	// 数据目录
	dataDir = strings.TrimSpace(dataDir)
	if !util_os.Stat(dataDir).IsExist() {
		util_os.MkDir(dataDir, os.ModePerm)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	arg = typ.Arg{
		TimeZone: timeZone,
		Port:     port,
		Path:     path,
		DataDir:  dataDir,
		AllowReg: strings.TrimSpace(allowReg) == "true",
	}
	typ.SetArg(arg)
	log.Println(arg)
}
