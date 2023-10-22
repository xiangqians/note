// 用户注销
// @author xiangqian
// @date 21:41 2023/07/11
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/context"
	"note/src/session"
)

// SignOut 注销
func SignOut(ctx *gin.Context) {
	// 获取当前用户
	user, _ := session.GetUser(ctx)

	// 清除session
	session.Clear(ctx)

	// 重定向
	session.Set(ctx, signInNameKey, user.Name)
	context.Redirect(ctx, "/user/signin", nil)
}
