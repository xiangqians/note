// arg
// @author xiangqian
// @date 00:11 2022/12/31
package app

import (
	"flag"
	"log"
	"note/src/api"
	"note/src/typ"
	"note/src/util"
	"path/filepath"
	"strings"
)

var arg typ.Arg

// 解析应用参数
func parseAppArg() {
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
	if !util.IsExistOfPath(dataDir) {
		util.Mkdir(dataDir)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	log.Printf("Port: %v\n", port)
	log.Printf("DataDir: %v\n", dataDir)
	log.Printf("AllowReg: %v\n", allowReg)

	arg = typ.Arg{
		Port:     port,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}

	// 设置api arg
	api.SetArg(arg)
}
