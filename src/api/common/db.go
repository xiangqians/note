// db
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_db "note/src/db"
	typ_page "note/src/typ/page"
	util_os "note/src/util/os"
)

func DbPage[T any](context *gin.Context, req typ_page.PageReq, sql string, args ...any) (typ_page.Page[T], error) {
	var page typ_page.Page[T]

	// db
	db := Db(context)
	defer db.Close()

	// open
	err := db.Open()
	if err != nil {
		return page, err
	}

	// begin
	err = db.Begin()
	if err != nil {
		return page, err
	}
	defer db.Commit()

	// page
	return _db.Page[T](db, req, sql, args...)
}

func DbQry[T any](context *gin.Context, sql string, args ...any) (T, int64, error) {
	var t T

	// db
	db := Db(context)
	defer db.Close()

	// open
	err := db.Open()
	if err != nil {
		return t, 0, err
	}

	// begin
	err = db.Begin()
	if err != nil {
		return t, 0, err
	}
	defer db.Commit()

	// qry & mapper
	return _db.RowsMapper[T](db.Qry(sql, args...))
}

func DbUpd(context *gin.Context, sql string, args ...any) (int64, error) {
	// db
	db := Db(context)
	defer db.Close()

	// open
	err := db.Open()
	if err != nil {
		return 0, err
	}

	// begin
	err = db.Begin()
	if err != nil {
		return 0, err
	}
	defer db.Commit()

	// upd
	affect, err := db.Upd(sql, args...)

	return affect, err
}

func DbDel(context *gin.Context, sql string, args ...any) (int64, error) {
	// db
	db := Db(context)
	defer db.Close()

	// open
	err := db.Open()
	if err != nil {
		return 0, err
	}

	// begin
	err = db.Begin()
	if err != nil {
		return 0, err
	}
	defer db.Commit()

	// del
	affect, err := db.Del(sql, args...)

	return affect, err
}

func DbAdd(context *gin.Context, sql string, args ...any) (int64, error) {
	// db
	db := Db(context)
	defer db.Close()

	// open
	err := db.Open()
	if err != nil {
		return 0, err
	}

	// begin
	err = db.Begin()
	if err != nil {
		return 0, err
	}
	defer db.Commit()

	// add
	id, err := db.Add(sql, args...)

	return id, err
}

// Db 获取db
func Db(context *gin.Context) _db.Db {
	// db
	dataDir := DataDir(context)
	dsn := fmt.Sprintf("%s%s%s", dataDir, util_os.FileSeparator, "database.db")
	return _db.Get(dsn)
}
