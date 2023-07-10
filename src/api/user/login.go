// login
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"github.com/gin-gonic/gin"
	app_context "note/src/context"
	"note/src/session"
	"note/src/typ"
)

// Login 登录页面
func Login(context *gin.Context) {
	resp, _ := session.Get[typ.Resp[typ.User]](context, app_context.RespSessionKey, true)
	app_context.HtmlOk(context, "user/login", resp)
}

// Login0 登录
func Login0(context *gin.Context) {

}
