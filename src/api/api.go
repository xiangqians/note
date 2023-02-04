// api
// @author xiangqian
// @date 14:52 2023/02/04
package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Html(pContext *gin.Context, templateName string, h gin.H, err error) {
	message, _ := SessionV[string](pContext, "message", true)

	if h == nil {
		h = gin.H{}
	}

	if err != nil {
		if message != "" {
			message += ", "
		}
		message += err.Error()
	}

	_, r := h["user"]
	if !r {
		user, _ := SessionUser(pContext)
		h["user"] = user
	}

	// 没有消息就是最好的消息
	h["message"] = message

	pContext.HTML(http.StatusOK, templateName, h)
}

func Redirect(pContext *gin.Context, location string, message any, m map[string]any) {
	if message != nil {
		if v, r := message.(error); r {
			message = v.Error()
		}
		if m == nil {
			m = map[string]any{}
		}
		m["message"] = message
	}

	if m != nil {
		session := sessions.Default(pContext)
		for k, v := range m {
			session.Set(k, v)
		}
		session.Save()
	}

	pContext.Redirect(http.StatusMovedPermanently, location)
}
