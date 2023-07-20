// sign out
// @author xiangqian
// @date 21:41 2023/07/11
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/config"
)

// SignOut 注销
func SignOut(context *gin.Context) {
	// user
	//user, _ := session.GetUser(context)

	// 清除session
	config.Clear(context)

	// 重定向
	context.Redirect(context, config.GetArg().Path+"/user/signin")
}
