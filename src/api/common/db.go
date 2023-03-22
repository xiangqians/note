// db
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/db"
	"note/src/typ"
	"note/src/util"
)

func DbQry[T any](context *gin.Context, sql string, args ...any) (T, int64, error) {
	return db.Qry[T](dsn(context), sql, args...)
}

func DbAdd(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Add(dsn(context), sql, args...)
}

func DbUpd(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Upd(dsn(context), sql, args...)
}

func DbDel(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Del(dsn(context), sql, args...)
}

func DbPage[T any](context *gin.Context, req typ.PageReq, sql string, args ...any) (typ.Page[T], error) {
	return db.Page[T](dsn(context), req, sql, args...)
}

func dsn(context *gin.Context) string {
	dataDir := DataDir(context)
	return fmt.Sprintf("%s%s%s", dataDir, util.FileSeparator, "database.db")
}
