// signup
// @author xiangqian
// @date 23:33 2023/07/10
package user

import (
	"github.com/gin-gonic/gin"
	app_context "note/src/context"
	"note/src/session"
	"note/src/typ"
)

// Signup 注册页
func Signup(context *gin.Context) {
	resp, _ := session.Get[typ.Resp[typ.User]](context, app_context.RespSessionKey, true)
	app_context.HtmlOk(context, "user/signup", resp)
}

// Signup0 注册
func Signup0(context *gin.Context) {

}
