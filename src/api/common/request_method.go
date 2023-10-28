// @author xiangqian
// @date 15:52 2023/10/28
package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/context"
)

// RequestMethod 获取请求方法
func RequestMethod(ctx *gin.Context) string {
	method, _ := context.Query[string](ctx, "_method")

	if method == "" {
		method, _ = context.PostForm[string](ctx, "_method")
	}

	switch method {
	case http.MethodGet:
		return http.MethodGet

	case http.MethodPost:
		return http.MethodPost

	case http.MethodPut:
		return http.MethodPut

	case http.MethodDelete:
		return http.MethodDelete

	default:
		return ctx.Request.Method
	}
}
