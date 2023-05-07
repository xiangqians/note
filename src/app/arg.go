// arg
// @author xiangqian
// @date 00:11 2022/12/31
package app

import (
	"flag"
	"log"
	api_common "note/src/api/common"
	"note/src/typ"
	util_os "note/src/util/os"
	"path/filepath"
	"strings"
)

var appArg typ.AppArg

// 解析应用参数
func parseArg() {
	var port int
	var dataDir string
	var allowReg int

	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"

	// parse
	flag.IntVar(&port, "port", 8080, "-port 8080")
	flag.StringVar(&dataDir, "dataDir", "./data", "-dataDir ./data")
	flag.IntVar(&allowReg, "allowReg", 1, "-allowReg 1")
	flag.Parse()

	// DataDir
	dataDir = strings.TrimSpace(dataDir)
	if !util_os.IsExist(dataDir) {
		util_os.MkDir(dataDir)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	log.Printf("Port: %v\n", port)
	log.Printf("DataDir: %v\n", dataDir)
	log.Printf("AllowReg: %v\n", allowReg)

	appArg = typ.AppArg{
		Port:     port,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}

	// 设置api App Arg
	api_common.AppArg = appArg
}
