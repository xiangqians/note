// arg
// @author xiangqian
// @date 00:11 2022/12/31
package app

import (
	"flag"
	"log"
	"note/src/util"
	"path/filepath"
	"strings"
)

var Port int       // 监听端口
var DataDir string // 数据目录
var AllowReg int   // 是否允许用户注册，0-不允许，1-允许

// 解析应用参数
func arg() {
	// parse
	flag.IntVar(&Port, "port", 8080, "-port 8080")
	flag.StringVar(&DataDir, "dataDir", "./data", "-dataDir ./data")
	flag.IntVar(&AllowReg, "allowReg", 1, "-allowReg 1")
	flag.Parse()

	// DataDir
	DataDir = strings.TrimSpace(DataDir)
	if !util.IsExistOfPath(DataDir) {
		util.Mkdir(DataDir)
	}
	// 获取绝对路径
	DataDir, _ = filepath.Abs(DataDir)

	log.Printf("Port: %v\n", Port)
	log.Printf("DataDir: %v\n", DataDir)
	log.Printf("AllowReg: %v\n", AllowReg)

	// -dataDir "C:\Users\xiangqian\Desktop\tmp\note\data"
}
