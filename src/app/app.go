// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common/trans"
)

// Run 运行Application
func Run() {
	// 日志记录器
	logger()

	// 解析应用参数
	arg()

	// 设置时区
	loc()

	// 检验器翻译
	trans.Translator()

	// gin模式
	gin.SetMode(gin.DebugMode)

	// default Engine
	engine := gin.Default()

	// template
	template(engine)

	// middleware
	middleware(engine)

	// route
	route(engine)

	// run
	addr := fmt.Sprintf(":%d", appArg.Port) // addr
	engine.Run(addr)
}
