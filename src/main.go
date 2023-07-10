// main
// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/arg"
	"note/src/local"
	"note/src/log"
	"note/src/middleware"
	"note/src/route"
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

	// 初始化中间件
	middleware.Init(engine)

	// route
	route.Init(engine)

	// run
	addr := fmt.Sprintf(":%d", arg.Arg.Port) // addr
	engine.Run(addr)
}
