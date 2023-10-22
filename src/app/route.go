// 路由
// @author xiangqian
// @date 18:11 2023/08/06
package app

import (
	"github.com/gin-gonic/gin"
	"note/src/api/image"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/context"
	"note/src/typ"
)

// 初始化路由
func initRoute(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(ctx *gin.Context) {
		context.HtmlNotFound(ctx, "404", typ.Resp[any]{})
	})

	contextPath := arg.ContextPath

	// user
	userGroup := engine.Group(contextPath + "/user")
	{
		userGroup.Any("/signin", user.SignIn)
		userGroup.Any("/signup", user.SignUp)
		userGroup.Any("/signout", user.SignOut)
		userGroup.Any("/settings", user.Settings)
	}

	// image
	imageGroup := engine.Group(contextPath + "/image")
	{
		imageGroup.Any("", image.List)
		imageGroup.Any("/rename", image.Rename)
	}

	// index
	engine.Any(contextPath+"/", index.Index)
}
