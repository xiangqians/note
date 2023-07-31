// signin
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/app"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
)

// SignIn 登录页
func SignIn(ctx *gin.Context) {
	context.HtmlOk(ctx, "user/signin", typ.Resp[typ.User]{})
}

// SignIn0 登录
func SignIn0(ctx *gin.Context) {

	app.ClearUser(1)

	// 保存用户信息到session
	session.SetUser(ctx, typ.User{Abs: typ.Abs{Id: 1}, Name: "test", Nickname: "测试"})

	// 重定向到首页
	context.Redirect(ctx, app.GetArg().Path+"/")
}
