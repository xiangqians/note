// auth
// @author xiangqian
// @date 23:19 2023/07/18
package config

import (
	"github.com/gin-gonic/gin"
	"net/http"
	context2 "note/src/context"
	"note/src/session"
	"strings"
)

// 初始化授权
func initAuth(engine *gin.Engine) {
	// 未授权拦截
	engine.Use(func(context *gin.Context) {
		// 请求路径
		reqPath := context.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, arg.Path+"/static/") {
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
		if reqPath == arg.Path+"/user/signin" || // 登录页
			reqPath == arg.Path+"/user/signup" || // 注册页
			(context.Request.Method == http.MethodPost && (reqPath == arg.Path+"/user/signin0" || reqPath == arg.Path+"/user/signup0")) { // 登录接口和注册接口
			// 如果已登录则重定向到首页
			if signIn {
				// 重定向到首页
				context2.Redirect(context, arg.Path+"/")
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
			context2.Redirect(context, arg.Path+"/user/signin")
			// 中止调用链
			context.Abort()
			return
		}
	})
}
