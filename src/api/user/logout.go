// user logout
// @author xiangqian
// @date 20:44 2023/05/11
package user

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/session"
	"note/src/typ"
)

// Logout 用户登出
func Logout(context *gin.Context) {
	// user
	user, _ := session.GetUser(context)

	// 清除session
	session.Clear(context)

	// 重定向
	api_common_context.Redirect(context, "/user/login", typ.Resp[typ.User]{
		Data: typ.User{Name: user.Name},
	})
}
