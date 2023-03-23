// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api"
	"note/src/api/common"
	api_img "note/src/api/img"
	api_index "note/src/api/index"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		common.Html(pContext, http.StatusNotFound, "404.html", nil, nil)
	}
	engine.Any("/404", handler)
	engine.NoRoute(handler)

	// index
	engine.Any("/", api_index.Page)

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
	imgGroup := engine.Group("/img")
	imgGroup.Any("/list", api_img.List) // page
	imgGroup.POST("/upload", api_img.Upload)
	imgGroup.PUT("/upload", api_img.Upload)
	imgGroup.PUT("/name", api_img.UpdName)
	imgGroup.DELETE("/:id", api_img.Del)
	imgGroup.GET("/:id", api_img.Get)
	imgGroup.Any("/:id/view", api_img.View) // page
	imgGroup.Any("/:id/edit", api_img.Edit) // page
}
