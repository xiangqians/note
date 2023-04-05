// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"github.com/gin-gonic/gin"
	api_common "note/src/api/common"
	api_img "note/src/api/img"
	api_index "note/src/api/index"
	api_note "note/src/api/note"
	api_recycle "note/src/api/recycle"
	api_user "note/src/api/user"
	typ_resp "note/src/typ/resp"
)

func route(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(context *gin.Context) {
		resp := typ_resp.Resp[any]{}
		api_common.HtmlNotFound(context, "404.html", resp)
	})

	// index
	engine.Any("/", api_index.Index)

	// user
	userGroup := engine.Group("/user")
	{
		userGroup.Any("/reg", api_user.Reg) // page
		userGroup.POST("", api_user.Add)
		userGroup.Any("/login", api_user.Login) // page
		userGroup.POST("/login0", api_user.Login0)
		userGroup.Any("/logout", api_user.Logout)
		userGroup.Any("/settings", api_user.Settings) // page
		userGroup.PUT("", api_user.Upd)
	}

	// note
	noteGroup := engine.Group("/note")
	{
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
	}

	// img
	imgGroup := engine.Group("/img")
	imgGroup.Any("/list", api_img.List) // page
	imgGroup.POST("/upload", api_img.Upload)
	imgGroup.PUT("/upload", api_img.Upload)
	imgGroup.GET("/:id", api_img.Get)
	imgGroup.GET("/:id/hist/:idx", api_img.GetHist)
	imgGroup.Any("/:id/view", api_img.View)               // page
	imgGroup.Any("/:id/hist/:idx/view", api_img.HistView) // page
	imgGroup.PUT("/name", api_img.UpdName)
	imgGroup.DELETE("/:id", api_img.Del)

	// recycle
	recycleGroup := engine.Group("/recycle")
	imgRecycleGroup := recycleGroup.Group("/img")
	imgRecycleGroup.Any("/list", api_recycle.ImgList) // page
	imgRecycleGroup.PUT("/:id/restore", api_recycle.ImgRestore)
	imgRecycleGroup.DELETE("/:id/permlyDel", api_recycle.ImgPermlyDel)

}
