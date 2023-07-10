// Route
// @author xiangqian
// @date 21:47 2022/12/23
package route

import (
	"github.com/gin-gonic/gin"
	"note/src/api/user"
	app_context "note/src/context"
	"note/src/typ"
)

// Init 初始化路由
func Init(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(context *gin.Context) {
		resp := typ.Resp[any]{}
		app_context.HtmlNotFound(context, "404.html", resp)
	})

	// user
	userGroup := engine.Group("/user")
	{
		userGroup.Any("/login", user.Login)
		userGroup.POST("/login0", user.Login0)
		//userGroup.Any("/signup", user.Signup)
		//userGroup.POST("/signup", user.Signup0)
		//userGroup.Any("/logout", user.Logout)
		//userGroup.Any("/settings", user.Settings) // page
		//userGroup.POST("/settings0", user.Settings0)
	}

	//// index
	//engine.Any("/", index.Index)

}
