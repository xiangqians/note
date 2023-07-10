// arg
// @author xiangqian
// @date 19:47 2023/07/10
package arg

import (
	"flag"
	"log"
	"note/src/typ"
	util_os "note/src/util/os"
	"path/filepath"
	"strings"
)

var Arg typ.Arg

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
	if !util_os.IsExist(dataDir) {
		util_os.MkDir(dataDir)
	}
	// 获取绝对路径
	dataDir, _ = filepath.Abs(dataDir)

	Arg = typ.Arg{
		Loc:      loc,
		Port:     port,
		DataDir:  dataDir,
		AllowReg: allowReg,
	}
	log.Println(Arg)
}
