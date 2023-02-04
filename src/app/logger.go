// logger
// @author xiangqian
// @date 19:32 2022/12/03
package app

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func Logger() {
	// 创建日志文件夹，如果不存在的话
	logDir := "./logs"
	fileInfo, err := os.Stat(logDir)
	if err != nil || !fileInfo.IsDir() {
		err = os.Mkdir(logDir, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 创建日志文件
	pLogFile, err := os.Create(logDir + "/debug.log")
	if err != nil {
		panic(err)
	}

	// gin mode
	gin.SetMode(gin.DebugMode)
	// 设置gin日志默认输出到：日志文件和控制台
	writer := io.MultiWriter(pLogFile, os.Stdout)
	gin.DefaultWriter = writer

	// log
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置日志输出
	log.SetOutput(writer)
}
