// @author xiangqian
// @date 23:05 2023/07/20
package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"note/src/model"
	util_string "note/src/util/string"
)

const UserKey = "__user__"

// GetUser 获取session用户信息
func GetUser(ctx *gin.Context) model.User {
	value := Get(ctx, UserKey, false)
	if user, ok := value.(model.User); ok {
		return user
	}
	return model.User{}
}

// SetUser 保存用户信息到session
func SetUser(ctx *gin.Context, user model.User) {
	Set(ctx, UserKey, user)
}

func SetString(ctx *gin.Context, name string, value any) {
	Set(ctx, name, util_string.String(value))
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

func GetString(ctx *gin.Context, name any, del bool) string {
	value := Get(ctx, name, del)
	return util_string.String(value)
}

// Get 根据名称获取session值
// name: 名称
// del : 是否删除session中的名称
func Get(ctx *gin.Context, name any, del bool) any {
	// 获取session
	session := sessions.Default(ctx)

	// 获取值
	value := session.Get(name)

	// 是否删除名称
	if del {
		// 删除名称
		session.Delete(name)
		// 保存，使删除生效
		session.Save()
	}

	return value
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
