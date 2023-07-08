// arg
// @author xiangqian
// @date 00:11 2022/12/31
package arg

import (
	"flag"
	"log"
	"note/src/typ"
	"note/src/util/os"
	"path/filepath"
	"strings"
)

var arg typ.Arg

// Init 初始化应用参数
func Init() {
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

	arg = typ.Arg{
		Loc:      loc,
		Port:     port,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}
	log.Println(arg)
}

func Get() typ.Arg {
	return arg
}
