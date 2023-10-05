// app
// @author xiangqian
// @date 22:39 2023/07/20
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Run 启动应用
func Run() {
	// 初始化日志记录器
	initLog()

	// 初始化应用参数
	initArg()

	// 初始化时区
	initTimeZone()

	// 设置gin模式
	gin.SetMode(gin.DebugMode)

	// 创建默认Engine
	engine := gin.Default()

	// 初始化模板
	initTemplate(engine)

	// 初始化session
	initSession(engine)

	// 初始化i18n
	initI18n(engine)

	// 初始化静态资源
	initStatic(engine)

	// 初始化授权
	initAuth(engine)

	// 初始化路由
	initRoute(engine)

	// 运行
	addr := fmt.Sprintf(":%d", arg.Port)
	engine.Run(addr)
}
