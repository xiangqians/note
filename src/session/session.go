// session
// @author xiangqian
// @date 23:05 2023/07/20
package session

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"note/src/typ"
)

const userKey = "__user__"

var Data *lru.Cache[string, map[any]any]

// GetUser 获取session用户信息
func GetUser(ctx *gin.Context) (typ.User, error) {
	user, err := Get[typ.User](ctx, userKey, false)

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user, err
}

// SetUser 保存用户信息到session
func SetUser(ctx *gin.Context, user typ.User) {
	// 单用户多端登录限制
	//data := Data
	//keys := data.Keys()
	//if keys != nil && len(keys) > 0 {
	//	log.Println(len(data.Values()), data.Values())
	//	for _, key := range keys {
	//		if value, r := data.Get(key); r {
	//			var v = value[userKey]
	//			if sessionUser, r := v.(typ.User); r {
	//				if sessionUser.Id == user.Id {
	//					data.Remove(key)
	//				}
	//			}
	//		}
	//	}
	//	log.Println(len(data.Values()), data.Values())
	//}
	Set(ctx, userKey, user)
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
	if t, r := value.(T); r {
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
