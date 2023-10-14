// Route
// @author xiangqian
// @date 21:47 2022/12/23
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	_db "note/src/db"
	"note/src/session"
	"note/src/typ"
	util_os "note/src/util/os"
)

// Db 获取数据库操作实例
func Db(ctx *gin.Context) (*gorm.DB, error) {
	dataDir := typ.GetArg().DataDir
	if ctx == nil {
		return _db.Db(util_os.Path(dataDir, "database.db"))
	}

	user, err := session.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	return _db.Db(util_os.Path(dataDir, fmt.Sprintf("%d", user.Id), "database.db"))
}

func DbExec(ctx *gin.Context, sql string, values ...any) (rowsAffected int64, err error) {
	db, err := Db(ctx)
	if err != nil {
		return
	}
	return _db.Exec(db, sql, values)
}

func DbRaw[T any](ctx *gin.Context, sql string, values ...any) (T, error) {
	db, err := Db(ctx)
	if err != nil {
		var t T
		return t, err
	}
	return _db.Raw[T](db, sql, values)
}

func DbPage[T any](ctx *gin.Context, current int64, size uint8, sql string, values ...any) (typ.Page[T], error) {
	db, err := Db(ctx)
	if err != nil {
		return typ.Page[T]{
			Current: current,
			Size:    size,
			Total:   0,
		}, err
	}
	return _db.Page[T](db, current, size, sql, values)
}
