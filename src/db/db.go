// DB
// https://pkg.go.dev/github.com/mattn/go-sqlite3
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	typ_page "note/src/typ"
	"strings"
)

// Page 分页查询
func Page[T any](dsn string, pageReq typ_page.Req, sql string, args ...any) (typ_page.Page[T], error) {
	// page
	current := pageReq.Current
	size := pageReq.Size
	page := typ_page.Page[T]{
		Current: current,
		Size:    size,
	}

	// db
	db := Get(dsn)
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

	// count
	countSql := fmt.Sprintf("SELECT COUNT(1) %s", sql[strings.Index(sql, "FROM"):])
	if strings.Contains(countSql, "GROUP BY") {
		countSql = fmt.Sprintf("SELECT COUNT(1) FROM (%s) r", countSql)
	}
	total, _, err := RowsMapper[int64](db.Qry(countSql, args...))
	if err != nil {
		return page, err
	}
	if total == 0 {
		return page, nil
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
	rows := size
	sql = fmt.Sprintf("%s LIMIT %v, %v", sql, offset, rows)

	// query
	data, count, err := RowsMapper[[]T](db.Qry(sql, args...))
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

func Qry[T any](dsn string, sql string, args ...any) (T, int64, error) {
	var t T

	// db
	db := Get(dsn)
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
	return RowsMapper[T](db.Qry(sql, args...))
}

func Upd(dsn string, sql string, args ...any) (int64, error) {
	return exec(dsn, func(db Db, sql string, args ...any) (int64, error) {
		return db.Upd(sql, args...)
	}, sql, args...)
}

func Del(dsn string, sql string, args ...any) (int64, error) {
	return exec(dsn, func(db Db, sql string, args ...any) (int64, error) {
		return db.Del(sql, args...)
	}, sql, args...)
}

func Add(dsn string, sql string, args ...any) (int64, error) {
	return exec(dsn, func(db Db, sql string, args ...any) (int64, error) {
		return db.Add(sql, args...)
	}, sql, args...)
}

func exec(dsn string, f func(db Db, sql string, args ...any) (int64, error), sql string, args ...any) (int64, error) {
	// db
	db := Get(dsn)
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

	// func
	return f(db, sql, args...)
}

// Get 获取db
// dsn: Data Source Name
func Get(dsn string) Db {
	return &DbImpl{
		driver: "sqlite3",
		dsn:    dsn,
	}
}

// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO

type Db interface {
	// Open 打开db
	Open() error

	// Begin 开启事务
	Begin() error

	// Add 新增
	// returns the integer generated by the database in response to a command.
	// Typically this will be from an "auto increment" column when inserting a new row.
	// Not all databases support this feature, and the syntax of such statements varies.
	// return insertId
	Add(sql string, args ...any) (int64, error)

	// Del 删除
	// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
	// return affect
	Del(sql string, args ...any) (int64, error)

	// Upd 更新
	// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
	Upd(sql string, args ...any) (int64, error)

	// Qry 查询
	Qry(sql string, args ...any) (*sql.Rows, error)

	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// Close 关闭资源
	Close() error
}
