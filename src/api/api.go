// Route
// @author xiangqian
// @date 21:47 2022/12/23
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"note/src/api/index"
	"note/src/api/user"
	"note/src/app"
	"note/src/context"
	"note/src/db"
	"note/src/session"
	"note/src/typ"
	util_os "note/src/util/os"
)

// Init 初始化API
func Init(engine *gin.Engine) {
	// 设置默认路由
	engine.NoRoute(func(ctx *gin.Context) {
		resp := typ.Resp[any]{}
		context.HtmlNotFound(ctx, "404", resp)
	})

	path := app.GetArg().Path

	// user
	userGroup := engine.Group(path + "/user")
	{
		userGroup.Any("/signin", user.SignIn)
		userGroup.POST("/signin0", user.SignIn0)
		userGroup.Any("/signup", user.SignUp)
		userGroup.POST("/signup0", user.SignUp0)
		userGroup.Any("/signout", user.SignOut)
		//userGroup.Any("/settings", user.Settings) // page
		//userGroup.POST("/settings0", user.Settings0)
	}

	// index
	engine.Any(path+"/", index.Index)
}

func Db(ctx *gin.Context) (*gorm.DB, error) {
	dataDir := app.GetArg().DataDir
	if ctx == nil {
		return db.Db(util_os.Path(dataDir, "database.db"))
	}

	user, err := session.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	return db.Db(util_os.Path(dataDir, fmt.Sprintf("%d", user.Id), "database.db"))
}

func Resp[T any](data T, anyMsg any) typ.Resp[T] {
	msg := ""
	if anyMsg != nil {
		if err, r := anyMsg.(error); r {
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("%v", anyMsg)
		}
	}

	return typ.Resp[T]{
		Data: data,
		Msg:  msg,
	}
}
