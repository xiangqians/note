// route
// @author xiangqian
// @date 18:11 2023/08/06
package app

import (
	"github.com/gin-gonic/gin"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/context"
	"note/src/typ"
)

// 初始化路由
func initRoute(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(ctx *gin.Context) {
		resp := typ.Resp[any]{}
		context.HtmlNotFound(ctx, "404", resp)
	})

	path := arg.Path

	// user
	userGroup := engine.Group(path + "/user")
	{
		userGroup.Any("/signin", user.SignIn)
		userGroup.POST("/signin0", user.SignIn0)
		userGroup.Any("/signup", user.SignUp)
		userGroup.POST("/signup0", user.SignUp0)
		userGroup.Any("/signout", user.SignOut)
		//userGroup.Any("/settings", user.Settings) // page
		//userGroup.POST("/settings0", user.Settings0)
	}

	// index
	engine.Any(path+"/", index.Index)
}
