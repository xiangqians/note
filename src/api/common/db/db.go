// db
// @author xiangqian
// @date 19:59 2023/03/22
package db

import (
	_sql "database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/db"
	"note/src/typ"
	util_os "note/src/util/os"
	"strings"
)

type dbExecType int8

const (
	add dbExecType = iota
	del
	upd
	qry
)

// DbPage 分页查询
// current: 当前页
// size: 页数量
func DbPage[T any](context *gin.Context, current int64, size uint8, sql string, args ...any) (typ.Page[T], error) {
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
	total, _, err := db.RowsMapper[int64](rows)
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
	data, count, err := db.RowsMapper[[]T](rows)
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

func DbQry[T any](context *gin.Context, sql string, args ...any) (T, int64, error) {
	// query
	rows, _, err := dbExec(context, qry, sql, args...)
	if err != nil {
		var t T
		return t, 0, err
	}

	// mapper
	return db.RowsMapper[T](rows)
}

func DbUpd(context *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	_, rowsAffected, err = dbExec(context, upd, sql, args...)
	return
}

func DbDel(context *gin.Context, sql string, args ...any) (rowsAffected int64, err error) {
	_, rowsAffected, err = dbExec(context, del, sql, args...)
	return
}

func DbAdd(context *gin.Context, sql string, args ...any) (lastInsertId int64, err error) {
	_, lastInsertId, err = dbExec(context, add, sql, args...)
	return
}

func dbExec(context *gin.Context, typ dbExecType, sql string, args ...any) (*_sql.Rows, int64, error) {
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

	var rows *_sql.Rows
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

func Db(context *gin.Context) (db.Db, error) {
	dataDir := common.DataDir(context)
	dsn := fmt.Sprintf("%s%s%s", dataDir, util_os.FileSeparator(), "database.db")
	return db.Get(dsn)
}
