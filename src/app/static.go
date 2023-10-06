// static
// @author xiangqian
// @date 23:17 2023/07/18
package app

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// 初始化静态资源
// https://github.com/gin-contrib/static
func initStatic(engine *gin.Engine) {
	engine.Use(static.Serve(arg.ContextPath+"/static", static.LocalFile("./res/static", false)))
}
