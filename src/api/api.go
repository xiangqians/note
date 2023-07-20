// Route
// @author xiangqian
// @date 21:47 2022/12/23
package api

import (
	"github.com/gin-gonic/gin"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/config"
	"note/src/context"
	"note/src/typ"
)

// Init 初始化API
func Init(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(ctx *gin.Context) {
		resp := typ.Resp[any]{}
		context.HtmlNotFound(ctx, "404", resp)
	})

	path := config.GetArg().Path

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
