// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api"
)

func route(pEngine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		pContext.HTML(http.StatusOK, "404.html", gin.H{})
	}
	pEngine.Any("/404", handler)
	pEngine.NoRoute(handler)

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.Any("/regpage", api.UserRegPage)
		userRouterGroup.Any("/loginpage", api.UserLoginPage)
		//userRouterGroup.POST("/login", api.UserLogin)
		//userRouterGroup.POST("/logout", api.UserLogout)
		//userRouterGroup.Any("/settingpage", api.UserSettingPage)
	}
	//pEngine.POST("/user", api.UserAdd)
	//pEngine.PUT("/user", api.UserUpd)

}
