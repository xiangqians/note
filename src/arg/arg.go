// arg
// @author xiangqian
// @date 00:11 2022/12/31
package arg

import (
	"flag"
	"log"
	"note/src/util"
	"strings"
)

var Port int       // 监听端口
var DataDir string // 数据目录

func Parse() {
	// parse
	flag.IntVar(&Port, "port", 8080, "-port 8080")
	flag.StringVar(&DataDir, "dataDir", "./data", "-dataDir ./data")
	flag.Parse()

	// DataDir
	DataDir = strings.TrimSpace(DataDir)
	if !util.IsExistOfPath(DataDir) {
		util.Mkdir(DataDir)
	}
	// "path/filepath"
	// 获取绝对路径
	//DataDir, _ = filepath.Abs(DataDir)

	log.Printf("Port: %v\n", Port)
	log.Printf("DataDir: %v\n", DataDir)

	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"
}
