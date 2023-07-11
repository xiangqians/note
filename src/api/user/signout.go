// sign out
// @author xiangqian
// @date 21:41 2023/07/11
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/arg"
	src_context "note/src/context"
	"note/src/session"
)

// SignOut 注销
func SignOut(context *gin.Context) {
	// user
	//user, _ := session.GetUser(context)

	// 清除session
	session.Clear(context)

	// 重定向
	src_context.Redirect(context, arg.Arg.Path+"/user/signin")
}
