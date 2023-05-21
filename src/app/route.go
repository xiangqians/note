// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/img"
	"note/src/api/index"
	"note/src/api/note"
	"note/src/api/user"
	"note/src/typ"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(context *gin.Context) {
		resp := typ.Resp[any]{}
		api_common_context.HtmlNotFound(context, "404.html", resp)
	})

	// user
	userGroup := engine.Group("/user")
	{
		userGroup.Any("/reg", user.Reg) // page
		userGroup.POST("/reg0", user.Reg0)
		userGroup.Any("/login", user.Login) // page
		userGroup.POST("/login0", user.Login0)
		userGroup.Any("/logout", user.Logout)
		userGroup.Any("/settings", user.Settings) // page
		userGroup.POST("/settings0", user.Settings0)
	}

	// index
	engine.Any("/", index.Index)

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

	// note
	noteGroup := engine.Group("/note")
	noteGroup.Any("/list", note.List) // page
	noteGroup.POST("", note.Add)
	noteGroup.POST("/upload", note.Upload)
	noteGroup.POST("/reUpload", note.ReUpload)
	noteGroup.POST("/updName", note.UpdName)
	noteGroup.GET("/:id", note.Get)
	noteGroup.GET("/:id/hist/:idx", note.GetHist)
	noteGroup.Any("/:id/view", note.View)               // page
	noteGroup.Any("/:id/hist/:idx/view", note.HistView) // page
	noteGroup.Any("/:id/edit", note.Edit)               // page
	noteGroup.POST("/updContent", note.UpdContent)
	noteGroup.POST("/cut/:srcId/to/:dstId", note.Cut)
	noteGroup.DELETE("/:id", note.Del)
	noteGroup.POST("/:id/restore", note.Restore)
	//noteGroup.DELETE("/:id/permlyDel", note.PermlyDel)

}
