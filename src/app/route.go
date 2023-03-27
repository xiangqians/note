// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api"
	"note/src/api/common"
	api_file "note/src/api/file"
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
	fileGroup.Any("/list", api_file.List) // page
	fileGroup.POST("", api_file.Add)
	fileGroup.POST("/upload", api_file.Upload)
	fileGroup.PUT("/upload", api_file.Upload)
	fileGroup.GET("/:id", api_file.Get)
	fileGroup.Any("/:id/view", api_file.View) // page
	fileGroup.Any("/:id/edit", api_file.Edit) // page
	fileGroup.PUT("/name", api_file.UpdName)
	fileGroup.PUT("/content", api_file.UpdContent)
	fileGroup.PUT("/cut/:srcId/to/:dstId", api_file.Cut)
	fileGroup.DELETE("/:id", api_file.Del)

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
