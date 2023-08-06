// Route
// @author xiangqian
// @date 21:47 2022/12/23
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"note/src/db"
	"note/src/session"
	"note/src/typ"
	util_os "note/src/util/os"
)

var Arg typ.Arg

func Db(ctx *gin.Context) (*gorm.DB, error) {
	dataDir := Arg.DataDir
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
