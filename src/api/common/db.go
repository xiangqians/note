// db
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/db"
	typ_page "note/src/typ/page"
	util_os "note/src/util/os"
)

func DbPage[T any](context *gin.Context, req typ_page.Req, sql string, args ...any) (typ_page.Page[T], error) {
	return db.Page[T](Dsn(context), req, sql, args...)
}

func DbQry[T any](context *gin.Context, sql string, args ...any) (T, int64, error) {
	return db.Qry[T](Dsn(context), sql, args...)
}

func DbUpd(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Upd(Dsn(context), sql, args...)
}

func DbDel(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Del(Dsn(context), sql, args...)
}

func DbAdd(context *gin.Context, sql string, args ...any) (int64, error) {
	return db.Add(Dsn(context), sql, args...)
}

// Dsn Data Source Name
func Dsn(context *gin.Context) string {
	dataDir := DataDir(context)
	return fmt.Sprintf("%s%s%s", dataDir, util_os.FileSeparator(), "database.db")
}
