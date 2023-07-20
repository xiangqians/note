// signup
// @author xiangqian
// @date 23:33 2023/07/10
package user

import (
	"github.com/gin-gonic/gin"
	"note/src/context"
	"note/src/typ"
)

// SignUp 注册页
func SignUp(ctx *gin.Context) {
	context.HtmlOk(ctx, "user/signup", typ.Resp[typ.User]{})
}

// SignUp0 注册
func SignUp0(context *gin.Context) {

}
