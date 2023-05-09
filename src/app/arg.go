// arg
// @author xiangqian
// @date 00:11 2022/12/31
package app

import (
	"flag"
	"log"
	api_common "note/src/api/common"
	"note/src/typ"
	"note/src/util/os"
	"path/filepath"
	"strings"
)

var appArg typ.AppArg

// 解析应用参数
func arg() {
	var loc string
	var port int
	var dataDir string
	var allowReg int

	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"

	// parse
	flag.StringVar(&loc, "loc", "Asia/Shanghai", "-loc Asia/Shanghai")
	flag.IntVar(&port, "port", 8080, "-port 8080")
	flag.StringVar(&dataDir, "dataDir", "./data", "-dataDir ./data")
	flag.IntVar(&allowReg, "allowReg", 1, "-allowReg 1")
	flag.Parse()

	// loc
	loc = strings.TrimSpace(loc)

	// DataDir
	dataDir = strings.TrimSpace(dataDir)
	if !os.IsExist(dataDir) {
		os.MkDir(dataDir)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	appArg = typ.AppArg{
		Loc:      loc,
		Port:     port,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}
	log.Println(appArg)

	// 设置api App Arg
	api_common.AppArg = appArg
}
