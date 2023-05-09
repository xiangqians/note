// db
// @author xiangqian
// @date 19:59 2023/03/22
package db

import (
	database_sql "database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	_db "note/src/db"
	"note/src/typ"
	"note/src/util/os"
	"strings"
)

type Type int8

const (
	add Type = iota
	del
	upd
	qry
)

// Page 分页查询
// current: 当前页
// size: 页数量
func Page[T any](context *gin.Context, current int64, size uint8, sql string, args ...any) (typ.Page[T], error) {
	// page
	page := typ.Page[T]{
		Current: current,
		Size:    size,
	}

	// db
	db, err := Db(context)
	if err != nil {
		return page, err
	}
	defer db.Close()

	// begin
	err = db.Begin()
	if err != nil {
		return page, err
	}
	defer db.Commit()

	// count
	countSql := fmt.Sprintf("SELECT COUNT(1) %s", sql[strings.Index(sql, "FROM"):])
	if strings.Contains(countSql, "GROUP BY") {
		countSql = fmt.Sprintf("SELECT COUNT(1) FROM (%s) r", countSql)
	}
	rows, err := db.Qry(countSql, args...)
	if err != nil {
		return page, err
	}
	total, _, err := _db.RowsMapper[int64](rows)
	if err != nil || total == 0 {
		return page, err
	}

	// set total & pages
	page.Total = total
	pages := total / int64(size)
	if total%int64(size) != 0 {
		pages += 1
	}
	page.Pages = pages

	// [offset,] rows
	offset := (current - 1) * int64(size)
	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, size)

	// data
	rows, err = db.Qry(sql, args...)
	if err != nil {
		return page, err
	}
	data, count, err := _db.RowsMapper[[]T](rows)
	if err != nil {
		return page, err
	}
	if count > 0 {
		// 不赋予指针数据，以访发生逃逸
		//page.Data = &data
		page.Data = data
	}

	return page, nil
}

func Qry[T any](context *gin.Context, sql string, args ...any) (T, int64, error) {
	// query
	rows, _, err := exec(context, qry, sql, args...)
	if err != nil {
		var t T
		return t, 0, err
	}

	// mapper
	return _db.RowsMapper[T](rows)
}

func Upd(context *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	_, rowsAffected, err = exec(context, upd, sql, args...)
	return
}

func Del(context *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	_, rowsAffected, err = exec(context, del, sql, args...)
	return
}

func Add(context *gin.Context, sql string, args ...any) (lastInsertId int64, err error) {
	_, lastInsertId, err = exec(context, add, sql, args...)
	return
}

func exec(context *gin.Context, typ Type, sql string, args ...any) (*database_sql.Rows, int64, error) {
	// db
	db, err := Db(context)
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	// begin
	err = db.Begin()
	if err != nil {
		return nil, 0, err
	}

	// exec
	var rows *database_sql.Rows
	var i int64 = 0
	switch typ {
	// add
	case add:
		i, err = db.Add(sql, args...)

	// del
	case del:
		i, err = db.Del(sql, args...)

	// upd
	case upd:
		i, err = db.Upd(sql, args...)

	// qry
	case qry:
		rows, err = db.Qry(sql, args...)
	}

	// Commit or rollback ?
	if err != nil {
		defer db.Rollback()
	} else {
		defer db.Commit()
	}

	return rows, i, err
}

func Db(context *gin.Context) (_db.Db, error) {
	dataDir := common.DataDir(context)
	dsn := fmt.Sprintf("%s%s%s", dataDir, os.FileSeparator(), "database.db")
	return _db.Get(dsn)
}
