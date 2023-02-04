// Session
// @author xiangqian
// @date 14:17 2023/02/04
package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionV 根据 key 获取 session value
// key: key
// del: 是否删除 session 中的key
func SessionV[T any](pContext *gin.Context, key any, del bool) (T, error) {
	session := sessions.Default(pContext)
	value := session.Get(key)
	if del {
		session.Delete(key)
		session.Save()
	}

	// t
	if t, r := value.(T); r {
		return t, nil
	}

	// default
	var t T
	return t, errors.New("unknown")
}

func SessionKv(pContext *gin.Context, key, value any) {
	session := sessions.Default(pContext)
	session.Set(key, value)
	session.Save()
}

func SessionClear(pContext *gin.Context) {
	// 解析session
	session := sessions.Default(pContext)
	// 清除session
	session.Clear()
	// 保存session数据
	session.Save()
}
