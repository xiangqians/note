// 授权
// @author xiangqian
// @date 23:19 2023/07/18
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/session"
	"note/src/util/time"
	"strings"
)

// 初始化授权
func initAuth(engine *gin.Engine) {
	// 根路径
	path := arg.Path

	// 未授权拦截
	engine.Use(func(ctx *gin.Context) {
		// 请求路径
		reqPath := ctx.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, path+"/static/") {
			ctx.Next()
			return
		}

		// 是否已登录
		isSignIn := false
		user, err := session.GetUser(ctx)
		if err == nil && user.Id > 0 {
			isSignIn = true
		}

		// 用户登录和注册放行
		if reqPath == path+"/user/signin" || // 登录页、登录接口
			reqPath == path+"/user/signup" { // 注册页、注册接口
			// 如果已登录则重定向到首页
			if isSignIn {
				// 重定向到首页
				redirect(ctx, path+"/")
				// 中止调用链
				ctx.Abort()
			} else
			// 如果未登录，放行登录或注册
			{
				ctx.Next()
			}
			return
		}

		// 未登录
		if !isSignIn {
			// 重定向到登录页
			//ctx.Request.URL.Path = path + "/user/signin"
			//engine.HandleContext(ctx)
			// OR
			redirect(ctx, path+"/user/signin")
			// 中止调用链
			ctx.Abort()
			return
		}
	})
}

func redirect(ctx *gin.Context, location string) {
	if strings.Contains(location, "?") {
		location = fmt.Sprintf("%s&t=%d", location, time.NowUnix())
	} else {
		location = fmt.Sprintf("%s?t=%d", location, time.NowUnix())
	}
	ctx.Redirect(http.StatusMovedPermanently, location)
}
