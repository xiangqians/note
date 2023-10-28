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
	"reflect"
	"strings"
)

func Page[T any](ctx *gin.Context, current int64, size uint8, sql string, values ...any) (model.Page[T], error) {
	db, err := Db(ctx)
	if err != nil {
		return model.Page[T]{
			Current: current,
			Size:    size,
			Total:   0,
		}, err
	}
	return _db.Page[T](db, current, size, sql, values...)
}

// GetPermlyDelId 查询永久删除的数据表id，以复用
func GetPermlyDelId[T any](ctx *gin.Context) (int64, error) {
	// 获取泛型类型
	var t T
	reflectType := reflect.TypeOf(t)

	// 结构体名称（此处即数据表名）
	name := strings.ToLower(reflectType.Name())

	// 查询
	return Get[int64](ctx, fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 2 LIMIT 1", name))
}

func Get[T any](ctx *gin.Context, sql string, values ...any) (T, error) {
	db, err := Db(ctx)
	if err != nil {
		var t T
		return t, err
	}
	return _db.Raw[T](db, sql, values...)
}

//func Upd(ctx *gin.Context, sql string, values ...any) (rowsAffected int64, err error) {
//	return Exec(ctx, sql, values...)
//}
//
//func Del(ctx *gin.Context, sql string, values ...any) (rowsAffected int64, err error) {
//	return Exec(ctx, sql, values...)
//}
//
//func Add(ctx *gin.Context, sql string, values ...any) (rowsAffected int64, err error) {
//	return Exec(ctx, sql, values...)
//}

func Exec(ctx *gin.Context, sql string, values ...any) (rowsAffected int64, lastInsertId int64, err error) {
	db, err := Db(ctx)
	if err != nil {
		return
	}
	return _db.Exec(db, sql, values...)
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
