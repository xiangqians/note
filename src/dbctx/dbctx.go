// @author xiangqian
// @date 23:11 2023/10/23
package dbctx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	_db "note/src/db"
	"note/src/model"
	"note/src/session"
	util_os "note/src/util/os"
)

func Page[T any](ctx *gin.Context, current int64, size uint8, sql string, args ...any) (model.Page[T], error) {
	db, err := Db(ctx)
	if err != nil {
		return model.Page[T]{
			Current: current,
			Size:    size,
			Total:   0,
		}, err
	}
	return _db.Page[T](db, current, size, sql, args...)
}

// GetPermlyDelId 获取永久删除的数据表id，以复用
// table : 数据表名
func GetPermlyDelId(ctx *gin.Context, table string) (int64, error) {
	return Get[int64](ctx, fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 2 LIMIT 1", table))
}

func Get[T any](ctx *gin.Context, sql string, args ...any) (T, error) {
	db, err := Db(ctx)
	if err != nil {
		var t T
		return t, err
	}
	return _db.Raw[T](db, sql, args...)
}

func Upd(ctx *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	rowsAffected, _, err = exec(ctx, sql, args...)
	return
}

func Del(ctx *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	rowsAffected, _, err = exec(ctx, sql, args...)
	return
}

func Add(ctx *gin.Context, sql string, args ...any) (lastInsertId int64, err error) {
	_, lastInsertId, err = exec(ctx, sql, args...)
	return
}

func exec(ctx *gin.Context, sql string, args ...any) (rowsAffected int64, lastInsertId int64, err error) {
	db, err := Db(ctx)
	if err != nil {
		return
	}
	return _db.Exec(db, sql, args...)
}

// Db 获取数据库操作实例
func Db(ctx *gin.Context) (*gorm.DB, error) {
	dataDir := model.GetArg().DataDir
	if ctx == nil {
		return _db.Db(util_os.Path(dataDir, "database.db"))
	}

	user, err := session.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	dsn := util_os.Path(dataDir, fmt.Sprintf("%d", user.Id), "database.db")
	return _db.Db(dsn)
}
