// sign out
// @author xiangqian
// @date 21:41 2023/07/11
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/config"
	"note/src/context"
	"note/src/session"
)

// SignOut 注销
func SignOut(ctx *gin.Context) {
	// user
	//user, _ := session.GetUser(context)

	// 清除session
	session.Clear(ctx)

	// 重定向
	context.Redirect(ctx, config.GetArg().Path+"/user/signin")
}
