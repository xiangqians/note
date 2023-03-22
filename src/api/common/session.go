// session
// @author xiangqian
// @date 20:01 2023/03/22
package common

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"note/src/typ"
)

const sessionKeyUser = "_user_"

// GetSessionUser 获取session用户信息
func GetSessionUser(context *gin.Context) (typ.User, error) {
	user, err := GetSessionV[typ.User](context, sessionKeyUser, false)

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user, err
}

// SetSessionUser 保存用户信息到session
func SetSessionUser(context *gin.Context, user typ.User) {
	SetSessionKv(context, sessionKeyUser, user)
}

// SetSessionKv 设置session kv
func SetSessionKv(context *gin.Context, key string, value any) {
	session := Session(context)
	session.Set(key, value)
	session.Save()
}

// GetSessionV 根据key获取session value
// key: key
// del: 是否删除session中的key
func GetSessionV[T any](context *gin.Context, key any, del bool) (T, error) {
	session := Session(context)
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

// ClearSession 清空session
func ClearSession(context *gin.Context) {
	// 解析session
	session := Session(context)
	// 清除session
	session.Clear()
	// 保存session数据
	session.Save()
}

func Session(context *gin.Context) sessions.Session {
	return sessions.Default(context)
}
