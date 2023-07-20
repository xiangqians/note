// config
// @author xiangqian
// @date 22:39 2023/07/20
package config

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化配置
func Init() {
	// 初始化日志记录器
	initLog()

	// 初始化应用参数
	initArg()

	// 初始化时区
	initLocal()
}

// InitEngine 初始化Engine配置
func InitEngine(engine *gin.Engine) {
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
}
