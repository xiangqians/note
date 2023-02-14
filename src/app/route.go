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
	pEngine.POST("/file/upload", api.FileUpload)
	pEngine.POST("/file/reupload", api.FileReUpload)
	pEngine.Any("/file/:id/view", api.FileView)
	pEngine.Any("/file/:id/viewpage", api.FileViewPage)
	pEngine.Any("/file/:id/editpage", api.FileEditPage)
	pEngine.PUT("/file/name", api.FileUpdName)
	pEngine.PUT("/file/content", api.FileUpdContent)
	pEngine.PUT("/file/cut/:srcId/to/:dstId", api.FileCut)
	pEngine.DELETE("/file/:id", api.FileDel)

	// img
	pEngine.Any("/img/listpage", api.ImgListPage)
	pEngine.POST("/img/upload", api.ImgUpload)
	pEngine.POST("/img/reupload", api.ImgReUpload)
	pEngine.PUT("/img/name", api.ImgUpdName)
	pEngine.DELETE("/img/:id", api.ImgDel)
	pEngine.Any("/img/:id/view", api.ImgView)
	pEngine.Any("/img/:id/viewpage", api.ImgViewPage)
	pEngine.Any("/img/:id/editpage", api.ImgEditPage)
}
