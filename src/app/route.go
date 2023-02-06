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
		//userRouterGroup.POST("/logout", api.UserLogout)
		//userRouterGroup.Any("/settingpage", api.UserSettingPage)
	}
	pEngine.POST("/user", api.UserAdd)
	//pEngine.PUT("/user", api.UserUpd)

	// dir
	dirRouterGroup := pEngine.Group("/dir")
	{
		dirRouterGroup.Any("/list", api.DirListPage)
		//dirRouterGroup.Any("/list/*ids", api.DirListPage) // ids路由全部模糊匹配
		dirRouterGroup.Any("/list/:pid", api.DirListPage)
	}

	// file
	pEngine.Any("/file/listpage", api.FileListPage)

}
