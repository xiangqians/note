// 会话
// @author xiangqian
// @date 23:05 2023/07/20
package session

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"note/src/model"
)

const UserKey = "__user__"

// GetUser 获取session用户信息
func GetUser(ctx *gin.Context) (model.User, error) {
	user, err := Get[model.User](ctx, UserKey, false)

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user, err
}

// SetUser 保存用户信息到session
func SetUser(ctx *gin.Context, user model.User) {
	Set(ctx, UserKey, user)
}

// Set 设置session <k, v>
func Set(ctx *gin.Context, key string, value any) {
	session := Session(ctx)
	session.Set(key, value)
	session.Save()
}

// Get 根据key获取session value
// key: key
// del: 是否删除session中的key
func Get[T any](ctx *gin.Context, key any, del bool) (T, error) {
	session := Session(ctx)
	value := session.Get(key)
	if del {
		session.Delete(key)
		session.Save()
	}

	// t
	if t, ok := value.(T); ok {
		return t, nil
	}

	// err
	var t T
	return t, errors.New(fmt.Sprintf("Type conversion error, %s", key))
}

// Clear 清空session
func Clear(ctx *gin.Context) {
	// 解析session
	session := Session(ctx)

	// 清除session
	session.Clear()

	// 保存session数据
	session.Save()
}

// Session 获取session
func Session(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}
