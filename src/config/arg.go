// arg
// @author xiangqian
// @date 19:47 2023/07/10
package config

import (
	"flag"
	"log"
	"note/src/typ"
	util_os "note/src/util/os"
	"path/filepath"
	"strings"
)

var arg typ.Arg

func GetArg() typ.Arg {
	return arg
}

// 初始化应用参数
func initArg() {
	var loc string
	var port int
	var path string
	var dataDir string
	var allowReg int

	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"

	// 解析参数
	flag.StringVar(&loc, "loc", "Asia/Shanghai", "-loc Asia/Shanghai")
	flag.IntVar(&port, "port", 8080, "-port 8080")
	flag.StringVar(&path, "path", "/", "-path /")
	flag.StringVar(&dataDir, "dataDir", "./data", "-dataDir ./data")
	flag.IntVar(&allowReg, "allowReg", 1, "-allowReg 1")
	flag.Parse()

	// 时区
	loc = strings.TrimSpace(loc)

	// 项目根路径
	path = strings.TrimSpace(path)
	if path == "/" {
		path = ""
	}

	// 数据目录
	dataDir = strings.TrimSpace(dataDir)
	if !util_os.IsExist(dataDir) {
		util_os.MkDir(dataDir)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	arg = typ.Arg{
		Loc:      loc,
		Port:     port,
		Path:     path,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}
	log.Println(arg)
}
