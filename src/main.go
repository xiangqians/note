// main
// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/arg"
	"note/src/auth"
	"note/src/i18n"
	"note/src/local"
	"note/src/log"
	"note/src/route"
	"note/src/session"
	"note/src/static"
	"note/src/template"
)

func main() {
	// 初始化日志记录器
	log.Init()

	// 初始化应用参数
	arg.Init()

	// 初始化时区
	local.Init()

	// gin模式
	gin.SetMode(gin.DebugMode)

	// default Engine
	engine := gin.Default()

	// 初始化模板
	template.Init(engine)

	// 初始化session
	session.Init(engine)

	// 初始化i18n
	i18n.Init(engine)

	// 初始化静态资源
	static.Init(engine)

	// 初始化授权
	auth.Init(engine)

	// 初始化路由
	route.Init(engine)

	// run
	addr := fmt.Sprintf(":%d", arg.Arg.Port)
	engine.Run(addr)
}
