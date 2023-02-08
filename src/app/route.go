// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"note/src/api"
)

func route(pEngine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		api.Html(pContext, "404.html", nil, nil)
	}
	pEngine.Any("/404", handler)
	pEngine.NoRoute(handler)

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.Any("/regpage", api.UserRegPage)
		userRouterGroup.Any("/loginpage", api.UserLoginPage)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.Any("/logout", api.UserLogout)
		userRouterGroup.Any("/stgpage", api.UserStgPage)
	}
	pEngine.POST("/user", api.UserAdd)
	pEngine.PUT("/user", api.UserUpd)

	// index
	pEngine.Any("/", api.IndexPage)

	// file
	pEngine.POST("/file", api.FileAdd)
	pEngine.Any("/file/:id/viewpage", api.FileViewPage)
	pEngine.Any("/file/:id/editpage", api.FileEditPage)
	pEngine.DELETE("/file/:id", api.FileDel)
	pEngine.PUT("/file/rename", api.FileRename)

}
