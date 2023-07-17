// signin
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/arg"
	src_context "note/src/context"
	"note/src/session"
	"note/src/typ"
)

// SignIn 登录页
func SignIn(context *gin.Context) {
	src_context.HtmlOk(context, "user/signin", typ.Resp[typ.User]{})
}

// SignIn0 登录
func SignIn0(context *gin.Context) {

	session.ClearUser(1)

	// 保存用户信息到session
	session.SetUser(context, typ.User{Abs: typ.Abs{Id: 1}, Name: "test", Nickname: "测试"})

	// 重定向到首页
	src_context.Redirect(context, arg.Arg.Path+"/")
}
