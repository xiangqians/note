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

	// 如果返回指针值，有可能会发生逃逸？
	//return &user

	return user, err
}

// SetUser 保存用户信息到session
func SetUser(ctx *gin.Context, user model.User) {
	Set(ctx, UserKey, user)
}

// Set 设置session <name, value>
func Set(ctx *gin.Context, name string, value any) {
	// 获取session
	session := sessions.Default(ctx)

	// 设置值
	session.Set(name, value)

	// 保存
	session.Save()
}

// Get 根据名称获取session值
// name: 名称
// del : 是否删除session中的名称
func Get[T any](ctx *gin.Context, name any, del bool) (T, error) {
	// 获取值
	session := sessions.Default(ctx)
	value := session.Get(name)

	// 是否删除名称
	if del {
		// 删除名称
		session.Delete(name)
		// 保存，使得删除生效
		session.Save()
	}

	// 类型转换
	if t, ok := value.(T); ok {
		return t, nil
	}

	var t T
	return t, errors.New(fmt.Sprintf("Type conversion error, %s", name))
}

// Clear 清空session数据
func Clear(ctx *gin.Context) {
	// 获取session
	session := sessions.Default(ctx)

	// 清除session数据
	session.Clear()

	// 保存session数据
	session.Save()
}
