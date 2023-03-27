// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api"
	"note/src/api/common"
	"note/src/api/file"
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
	fileGroup := engine.Group("/file")
	fileGroup.Any("/list", file.List) // page
	engine.POST("/file", file.FileAdd)
	engine.POST("/file/upload", file.FileUpload)
	engine.POST("/file/reupload", file.FileReUpload)
	engine.Any("/file/:id/view", file.FileView)
	engine.Any("/file/:id/viewpage", file.FileViewPage)
	engine.Any("/file/:id/editpage", file.FileEditPage)
	engine.PUT("/file/name", file.FileUpdName)
	engine.PUT("/file/content", file.FileUpdContent)
	engine.PUT("/file/cut/:srcId/to/:dstId", file.FileCut)
	engine.DELETE("/file/:id", file.FileDel)

	// img
	imgGroup := engine.Group("/img")
	imgGroup.Any("/list", api_img.List) // page
	imgGroup.POST("/upload", api_img.Upload)
	imgGroup.PUT("/upload", api_img.Upload)
	imgGroup.PUT("/name", api_img.UpdName)
	imgGroup.DELETE("/:id", api_img.Del)
	imgGroup.GET("/:id", api_img.Get)
	imgGroup.Any("/:id/view", api_img.View) // page
}
