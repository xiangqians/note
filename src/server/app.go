// app
// @author xiangqian
// @date 22:39 2023/07/20
package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/ini"
)

// Run 启动应用
func Run() {
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
	addr := fmt.Sprintf(":%d", ini.Ini.Server.Port)
	engine.Run(addr)
}
