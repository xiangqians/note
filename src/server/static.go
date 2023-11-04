// static
// @author xiangqian
// @date 23:17 2023/07/18
package server

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"note/src/ini"
)

// 初始化静态资源
// https://github.com/gin-contrib/static
func initStatic(engine *gin.Engine) {
	engine.Use(static.Serve(ini.Ini.Server.ContextPath+"/static", static.LocalFile("./res/static", false)))
}
