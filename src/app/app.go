// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api_common "note/src/api/common"
	"strconv"
)

func Run() {
	// 日志记录器
	logger()

	// 解析应用参数
	parseArg()

	// 设置时区
	local()

	// ValidateTrans
	api_common.ValidateTrans()

	// gin模式：DebugMode、TestMode、ReleaseMode
	gin.SetMode(gin.DebugMode)

	// default Engine
	engine := gin.Default()

	// template
	htmlTemplate(engine)

	// middleware
	sessionMiddleware(engine)
	i18nMiddleware(engine)
	staticMiddleware(engine)
	permMiddleware(engine)

	// route
	route(engine)

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(arg.Port), 10))

	// run
	engine.Run(addr)
}
