// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api"
	"note/src/api/common"
	api_index "note/src/api/index"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		common.Html(pContext, http.StatusNotFound, "404.html", nil, nil)
	}
	engine.Any("/404", handler)
	engine.NoRoute(handler)

	// user
	userRouterGroup := engine.Group("/user")
	{
		userRouterGroup.Any("/regpage", api.UserRegPage)
		userRouterGroup.Any("/loginpage", api.UserLoginPage)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.Any("/logout", api.UserLogout)
		userRouterGroup.Any("/stgpage", api.UserStgPage)
	}
	engine.POST("/user", api.UserAdd)
	engine.PUT("/user", api.UserUpd)

	// index
	engine.Any("/", api_index.IndexPage)

	// file
	engine.Any("/file/listpage", api.FileListPage)
	engine.POST("/file", api.FileAdd)
	engine.POST("/file/upload", api.FileUpload)
	engine.POST("/file/reupload", api.FileReUpload)
	engine.Any("/file/:id/view", api.FileView)
	engine.Any("/file/:id/viewpage", api.FileViewPage)
	engine.Any("/file/:id/editpage", api.FileEditPage)
	engine.PUT("/file/name", api.FileUpdName)
	engine.PUT("/file/content", api.FileUpdContent)
	engine.PUT("/file/cut/:srcId/to/:dstId", api.FileCut)
	engine.DELETE("/file/:id", api.FileDel)

	// img
	engine.Any("/img/listpage", api.ImgListPage)
	engine.POST("/img/upload", api.ImgUpload)
	engine.POST("/img/reupload", api.ImgReUpload)
	engine.PUT("/img/name", api.ImgUpdName)
	engine.DELETE("/img/:id", api.ImgDel)
	engine.Any("/img/:id/view", api.ImgView)
	engine.Any("/img/:id/viewpage", api.ImgViewPage)
	engine.Any("/img/:id/editpage", api.ImgEditPage)
}
