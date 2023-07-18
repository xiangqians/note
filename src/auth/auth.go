// auth
// @author xiangqian
// @date 23:19 2023/07/18
package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/arg"
	src_context "note/src/context"
	"note/src/session"
	"strings"
)

// Init 初始化授权
func Init(engine *gin.Engine) {
	// 未授权拦截
	engine.Use(func(context *gin.Context) {
		// 请求路径
		reqPath := context.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, arg.Arg.Path+"/static/") {
			context.Next()
			return
		}

		// 是否已登录
		signIn := false
		user, err := session.GetUser(context)
		if err == nil && user.Id > 0 {
			signIn = true
		}

		// 用户登录和注册放行
		if reqPath == arg.Arg.Path+"/user/signin" || // 登录页
			reqPath == arg.Arg.Path+"/user/signup" || // 注册页
			(context.Request.Method == http.MethodPost && (reqPath == arg.Arg.Path+"/user/signin0" || reqPath == arg.Arg.Path+"/user/signup0")) { // 登录接口和注册接口
			// 如果已登录则重定向到首页
			if signIn {
				// 重定向到首页
				src_context.Redirect(context, arg.Arg.Path+"/")
				// 中止调用链
				context.Abort()
			} else
			// 如果未登录，放行登录或注册
			{
				context.Next()
			}
			return
		}

		// 未登录
		if !signIn {
			// 重定向到登录页
			//context.Request.URL.Path = arg.Arg.Path+"/user/signin"
			//engine.HandleContext(context)
			// OR
			src_context.Redirect(context, arg.Arg.Path+"/user/signin")
			// 中止调用链
			context.Abort()
			return
		}
	})
}
