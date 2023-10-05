// 日志
// @author xiangqian
// @date 19:45 2023/07/10
package app

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	util_os "note/src/util/os"
	"os"
	"path/filepath"
)

// 初始化日志记录器
func initLog() {
	// 当前目录
	curDir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	// 创建日志目录（如果文件不存在或者不是目录文件的话）
	logDir := util_os.Path(curDir, "log")
	file := util_os.Stat(logDir)
	if !file.IsExist() || !file.IsDir() {
		err = util_os.MkDir(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建日志文件（如果存在则覆盖）
	logFile, err := os.Create(util_os.Path(logDir, "debug.log"))
	if err != nil {
		panic(err)
	}

	// 多重写入器：写到文件 & 写到控制台
	writer := io.MultiWriter(logFile, os.Stdout)

	// 设置gin日志默认输出到：日志文件和控制台
	gin.DefaultWriter = writer

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置日志输出
	log.SetOutput(writer)

	log.Printf("logDir  %s\n", logDir)
	log.Printf("logFile %s\n", logFile.Name())
}
