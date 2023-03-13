// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"strconv"
)

func Run() {
	// 日志记录器
	logger()

	// 设置时区
	local()

	// 解析应用参数
	arg()

	// ValidateTrans
	api.ValidateTrans()

	// gin模式：DebugMode、ReleaseMode、TestMode
	gin.SetMode(gin.DebugMode)

	// default Engine
	pEngine := gin.Default()

	// template
	htmlTemplate(pEngine)

	// middleware
	sessionMiddleware(pEngine)
	i18nMiddleware(pEngine)
	staticMiddleware(pEngine)
	permMiddleware(pEngine)

	// route
	route(pEngine)

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(Port), 10))

	// run
	pEngine.Run(addr)
}
