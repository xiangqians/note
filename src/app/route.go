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
	api_note "note/src/api/note"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		common.Html(pContext, http.StatusNotFound, "404.html", nil, nil)
	}
	engine.Any("/404", handler)
	engine.NoRoute(handler)

	// index
	engine.Any("/", api_index.Index)

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

	// note
	noteGroup := engine.Group("/note")
	noteGroup.Any("/list", api_note.List) // page
	noteGroup.POST("", api_note.Add)
	noteGroup.POST("/upload", api_note.Upload)
	noteGroup.PUT("/upload", api_note.Upload)
	noteGroup.GET("/:id", api_note.Get)
	noteGroup.Any("/:id/view", api_note.View) // page
	noteGroup.Any("/:id/edit", api_note.Edit) // page
	noteGroup.PUT("/name", api_note.UpdName)
	noteGroup.PUT("/content", api_note.UpdContent)
	noteGroup.PUT("/cut/:srcId/to/:dstId", api_note.Cut)
	noteGroup.DELETE("/:id", api_note.Del)

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
