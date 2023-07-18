// log
// @author xiangqian
// @date 19:45 2023/07/10
package log

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	util_os "note/src/util/os"
	"os"
	"path/filepath"
)

// Init 初始化日志记录器
func Init() {
	// current directory
	curDir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	// 创建日志文件夹，如果不存在的话
	logDir := curDir + util_os.FileSeparator() + "log"
	fileInfo, err := os.Stat(logDir)
	if err != nil || !fileInfo.IsDir() {
		err = os.Mkdir(logDir, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 创建日志文件（如果存在则覆盖）
	logFile, err := os.Create(logDir + util_os.FileSeparator() + "debug.log")
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

	log.Printf("logDir: %s\n", logDir)
}
