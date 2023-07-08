// log
// @author xiangqian
// @date 21:11 2023/06/12
package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	util_os "note/src/util/os"
	"os"
	"path/filepath"
)

// 日志记录器
func init() {
	// current directory
	curDir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	// 创建日志文件夹，如果不存在的话
	logDir := fmt.Sprintf("%s%s%s", curDir, util_os.FileSeparator(), "log")
	fileInfo, err := os.Stat(logDir)
	if err != nil || !fileInfo.IsDir() {
		err = os.Mkdir(logDir, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 创建日志文件（如果存在则覆盖）
	logFile, err := os.Create(logDir + "/debug.log")
	if err != nil {
		panic(err)
	}

	// 设置gin日志默认输出到：日志文件和控制台
	writer := io.MultiWriter(logFile, os.Stdout)
	gin.DefaultWriter = writer

	// log
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置日志输出
	log.SetOutput(writer)

	// print log dir
	log.Printf("logDir: %s\n", logDir)
}
