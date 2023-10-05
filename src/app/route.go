// 路由
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

	// 服务根路径
	path := arg.Path

	// user
	userGroup := engine.Group(path + "/user")
	{
		userGroup.Any("/signIn", user.SignIn)
		userGroup.Any("/signUp", user.SignUp)
		userGroup.Any("/signOut", user.SignOut)
		//userGroup.Any("/settings", user.Settings) // page
		//userGroup.POST("/settings0", user.Settings0)
	}

	// index
	engine.Any(path+"/", index.Index)
}
