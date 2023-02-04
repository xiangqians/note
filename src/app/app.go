// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/arg"
	"strconv"
)

func Run() {
	// Logger
	Logger()

	// arg parse
	arg.Parse()

	// Gin ReleaseMode
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

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
