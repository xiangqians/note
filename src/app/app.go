// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/arg"
	"strconv"
)

func Run() {
	// Logger
	Logger()

	// 设置时区
	Local()

	// parse arg
	arg.Parse()

	// ValidateTrans
	api.ValidateTrans()

	// Gin ReleaseMode
	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)

	// default Engine
	pEngine := gin.Default()

	// init
	htmlTemplate(pEngine)
	sessionMiddleware(pEngine)
	i18nMiddleware(pEngine)
	staticMiddleware(pEngine)
	permMiddleware(pEngine)
	route(pEngine)

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(arg.Port), 10))

	// run
	pEngine.Run(addr)
}
