// Route
// @author xiangqian
// @date 21:47 2022/12/23
package api

import (
	"github.com/gin-gonic/gin"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/config"
)

// Init 初始化路由
func Init(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(context *gin.Context) {
		resp := context.Resp[any]{}
		context.HtmlNotFound(context, "404", resp)
	})

	// user
	userGroup := engine.Group(config.Get().Path + "/user")
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
	engine.Any(Arg.Path+"/", index.Index)

}
