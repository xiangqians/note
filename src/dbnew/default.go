// @author xiangqian
// @date 17:43 2023/10/29
package dbnew

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"strings"
	"sync"
)

type DefaultDb struct {
	db *sql.DB // db
	tx *sql.Tx // 事务支持
}

func (db *DefaultDb) Begin() (err error) {
	db.tx, err = db.db.Begin()
	return
}

func (db *DefaultDb) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	log.Println("[?ms]", sql)
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	insertId, _ = result.LastInsertId()
	return
}

func (db *DefaultDb) Del(sql string, args ...any) (rowsAffected int64, err error) {
	log.Println("[?ms]", sql)
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return
}

func (db *DefaultDb) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	log.Println("[?ms]", sql)
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return
}

func (db *DefaultDb) exec(sql string, args ...any) (sql.Result, error) {
	if db.tx != nil {
		return db.tx.Exec(sql, args...)
	}
	return db.db.Exec(sql, args...)
}

func (db *DefaultDb) Get(sql string, args ...any) (Result, error) {
	log.Println("[?ms]", sql)
	rows, err := db.query(sql, args...)
	return &DefaultResult{rows: rows}, err
}

func (db *DefaultDb) query(sql string, args ...any) (*sql.Rows, error) {
	if db.tx != nil {
		return db.tx.Query(sql, args...)
	}
	return db.db.Query(sql, args...)
}

func (db *DefaultDb) Page(sql string, current int64, size uint8, args ...any) (Result, error) {
	log.Println("[?ms]", sql)

	// 计数
	index := strings.Index(sql, "FROM")
	result, err := db.Get(fmt.Sprintf("SELECT COUNT(1) %s", sql[index:]), args...)
	if err != nil {
		return &DefaultResult{}, err
	}
	var count int64
	result.Scan(&count)

	// 查询分页数据
	offset := (current - 1) * int64(size)
	limit := size
	result, err = db.Get(fmt.Sprintf("%s LIMIT %d,%d", sql, offset, limit), args...)

	defaultResult := result.(*DefaultResult)
	defaultResult.count = count

	return defaultResult, err
}

func (db *DefaultDb) Commit() error {
	if db.tx != nil {
		return db.tx.Commit()
	}
	return nil
}

func (db *DefaultDb) Rollback() error {
	if db.tx != nil {
		return db.tx.Rollback()
	}
	return nil
}

func (db *DefaultDb) Close() error {
	//if db.db != nil {
	//	err := db.db.Close()
	//	db.db = nil
	//	return err
	//}
	return nil
}

type DefaultResult struct {
	count int64
	rows  *sql.Rows
}

func (result DefaultResult) Count() int64 {
	return result.count
}

func (result DefaultResult) Scan(dest any) error {
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer result.rows.Close()

	destType := reflect.TypeOf(dest).Elem()
	kind := destType.Kind()

	switch kind {
	// 基本数据类型
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		for result.rows.Next() {
			return result.rows.Scan(dest)
		}

	// 基本数据类型切片

	// 结构体
	case reflect.Struct:
		// 查询字段名称集
		//cols, err := result.rows.Columns()
		//if err != nil {
		//	return err
		//}
		//
		//if result.rows.Next() {
		//}

	// 结构体切片

	// map[string]any
	case reflect.Map:
	}

	return nil
}

type DefaultDbConnPool struct {
	Driver string     // driver name
	Dsn    string     // data source name
	mutex  sync.Mutex // sync.Mutex 是一个基本的同步原语，可以实现并发环境下的线程安全
	slice  []*sql.DB  // db切片
}

func (dbConnPool *DefaultDbConnPool) Get() (Db, error) {
	dbConnPool.mutex.Lock() // 获取锁
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer dbConnPool.mutex.Unlock() // 释放锁

	if dbConnPool.slice == nil {
		// len 0, cap ?
		dbConnPool.slice = make([]*sql.DB, 0, 16)
	}

	if len(dbConnPool.slice) > 0 {
		return &DefaultDb{db: dbConnPool.slice[0]}, nil
	}

	db, err := dbConnPool.open()
	if err != nil {
		return nil, err
	}

	dbConnPool.slice = append(dbConnPool.slice, db)
	return &DefaultDb{db: db}, nil
}

// 打开数据库连接
func (dbConnPool *DefaultDbConnPool) open() (*sql.DB, error) {
	log.Println("open db", dbConnPool.Dsn)

	db, err := sql.Open(dbConnPool.Driver, dbConnPool.Dsn)
	if err != nil {
		return db, err
	}

	// 校验数据库连接
	err = db.Ping()

	return db, err
}
