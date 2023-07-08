// session
// @author xiangqian
// @date 20:01 2023/03/22
package session

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"note/src/typ"
)

const userSessionKey = "__user__"

// GetUser 获取session用户信息
func GetUser(context *gin.Context) (typ.User, error) {
	user, err := Get[typ.User](context, userSessionKey, false)

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user, err
}

// SetUser 保存用户信息到session
func SetUser(context *gin.Context, user typ.User) {
	Set(context, userSessionKey, user)
}

// Set 设置session <k, v>
func Set(context *gin.Context, key string, value any) {
	session := Session(context)
	session.Set(key, value)
	session.Save()
}

// Get 根据key获取session value
// key: key
// del: 是否删除session中的key
func Get[T any](context *gin.Context, key any, del bool) (T, error) {
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

	// err
	var t T
	return t, errors.New(fmt.Sprintf("Unknown %s", key))
}

// Clear 清空session
func Clear(context *gin.Context) {
	// 解析session
	session := Session(context)

	// 清除session
	session.Clear()

	// 保存session数据
	session.Save()
}

// Session 获取session
func Session(context *gin.Context) sessions.Session {
	return sessions.Default(context)
}
