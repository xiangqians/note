// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	note2 "note/app/api/note"
	api_common_context "note/src/api/common/context"
	"note/src/api/img"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/typ"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(context *gin.Context) {
		resp := typ.Resp[any]{}
		api_common_context.HtmlNotFound(context, "404.html", resp)
	})

	// index
	engine.Any("/", index.Index)

	// user
	userGroup := engine.Group("/user")
	{
		userGroup.Any("/reg", user.Reg) // page
		userGroup.POST("", user.Add)
		userGroup.Any("/login", user.Login) // page
		userGroup.POST("/login0", user.Login0)
		userGroup.Any("/logout", user.Logout)
		userGroup.Any("/settings", user.Settings) // page
		userGroup.PUT("", user.Upd)
	}

	// note
	noteGroup := engine.Group("/note")
	noteGroup.Any("/list", note2.List) // page
	noteGroup.POST("", note2.Add)
	noteGroup.POST("/upload", note2.Upload)
	noteGroup.POST("/reUpload", note2.ReUpload)
	noteGroup.PUT("/name", note2.UpdName)
	noteGroup.GET("/:id", note2.Get)
	noteGroup.Any("/:id/view", note2.View) // page
	noteGroup.PUT("/content", note2.UpdContent)
	noteGroup.PUT("/cut/:srcId/to/:dstId", note2.Cut)
	noteGroup.DELETE("/:id", note2.Del)
	noteGroup.PUT("/:id/restore", note2.Restore)

	// img
	imgGroup := engine.Group("/img")
	imgGroup.Any("/list", img.List) // page
	imgGroup.POST("/upload", img.Upload)
	imgGroup.POST("/reUpload", img.ReUpload)
	imgGroup.GET("/:id", img.Get)
	imgGroup.GET("/:id/hist/:idx", img.GetHist)
	imgGroup.Any("/:id/view", img.View)               // page
	imgGroup.Any("/:id/hist/:idx/view", img.HistView) // page
	imgGroup.PUT("/name", img.UpdName)
	imgGroup.DELETE("/:id", img.Del)
	imgGroup.PUT("/:id/restore", img.Restore)
	imgGroup.DELETE("/:id/permlyDel", img.PermlyDel)

}
