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

func (result *DefaultResult) Count() int64 {
	return result.count
}

func (result *DefaultResult) Scan(dest any) error {
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer result.rows.Close()

	t := reflect.TypeOf(dest).Elem()
	switch t.Kind() {
	// 基本数据类型
	case // 布尔型
		reflect.Bool,
		// 整型
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		// 浮点型
		reflect.Float32, reflect.Float64,
		// 字符串类型
		reflect.String:
		if result.next() {
			return result.scanDefault(dest)
		}

	// 结构体
	case reflect.Struct:
		return result.scanStruct(dest)

	// 切片
	case reflect.Slice:
		t := t.Elem()
		switch t.Kind() {
		// 基本数据类型切片
		case // 布尔型
			reflect.Bool,
			// 整型
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			// 浮点型
			reflect.Float32, reflect.Float64,
			// 字符串类型
			reflect.String:

		// 结构体切片
		case reflect.Struct:
			// len 0, cap ?
			slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 1)
			for result.next() {
				e := reflect.New(t)
				result.scanStruct(e.Interface())
				slice = reflect.Append(slice, e.Elem())
			}
			reflect.ValueOf(dest).Elem().Set(slice)
		}
	}

	return nil
}

func (result *DefaultResult) next() bool {
	return result.rows.Next()
}

// 扫描基本数据类型
func (result *DefaultResult) scanDefault(dest any) error {
	return result.rows.Scan(dest)
}

// 扫描结构体
func (result *DefaultResult) scanStruct(dest any) error {
	// 查询字段集
	var err error
	cols, err := result.rows.Columns()
	if err != nil {
		return err
	}

	length := len(cols)

	// len ?, cap ?
	newDest := make([]any, length, length)
	var n any
	for i := 0; i < length; i++ {
		cols[i] = strings.ReplaceAll(cols[i], "_", "")
		newDest[i] = &n
	}

	t := reflect.TypeOf(dest).Elem()
	v := reflect.ValueOf(dest).Elem()
	result.dest(&cols, &newDest, t, v)

	return result.rows.Scan(newDest...)
}

func (result *DefaultResult) dest(cols *[]string, dest *[]any, t reflect.Type, v reflect.Value) {
	for i, length := 0, t.NumField(); i < length; i++ {
		field := t.Field(i)
		switch field.Type.Kind() {
		// 基本数据类型
		case // 布尔型
			reflect.Bool,
			// 整型
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			// 浮点型
			reflect.Float32, reflect.Float64,
			// 字符串类型
			reflect.String:
			for i, length := 0, len(*cols); i < length; i++ {
				col := (*cols)[i]
				// 不区分大小写比较
				if strings.EqualFold(col, field.Name) {
					//v := v.Field(i)
					v := v.FieldByName(field.Name)
					if v.CanAddr() {
						(*dest)[i] = v.Addr().Interface()
					}
				}
			}

		// 结构体
		case reflect.Struct:
			v := v.FieldByName(field.Name)
			result.dest(cols, dest, field.Type, v)
		}
	}
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
