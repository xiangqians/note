// @author xiangqian
// @date 17:43 2023/10/29
package dbnew

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	util_crypto_md5 "note/src/util/crypto/md5"
	"sync"
)

type DefaultDb struct {
	db *sql.DB
}

func (db *DefaultDb) Begin() error {
	return nil
}

func (db *DefaultDb) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	return 0, 0, nil
}

func (db *DefaultDb) Del(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *DefaultDb) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *DefaultDb) Get(sql string, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *DefaultDb) Page(sql string, current int64, size uint8, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *DefaultDb) Commit() error {
	return nil
}

func (db *DefaultDb) Rollback() error {
	return nil
}

func (db *DefaultDb) Close() error {
	return nil
}

type DefaultResult struct {
	count int64
}

func (result DefaultResult) Count() int64 {
	return result.count
}

func (result DefaultResult) Scan(dest any) error {
	//if result.tx == nil {
	//	return nil
	//}
	//
	//// Scan? Take?
	//return result.tx.Scan(dest).Error
	return nil
}

type DefaultDbConnPool struct {
	// sync.Mutex 是一个基本的同步原语，可以实现并发环境下的线程安全
	mutex sync.Mutex

	// db map
	m map[string]*sql.DB
}

func (dbConnPool *DefaultDbConnPool) Get(dsn string) (Db, error) {
	dbConnPool.mutex.Lock() // 获取锁
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer dbConnPool.mutex.Unlock() // 释放锁

	if dbConnPool.m == nil {
		// len 0, cap ?
		// cap ?
		dbConnPool.m = make(map[string]*sql.DB, 16)
	}

	key := util_crypto_md5.Encrypt([]byte(dsn), nil)
	if db, ok := dbConnPool.m[key]; ok {
		return &DefaultDb{db: db}, nil
	}

	db, err := dbConnPool.open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	dbConnPool.m[key] = db
	return &DefaultDb{db: db}, nil
}

// 打开数据库连接
func (dbConnPool *DefaultDbConnPool) open(driver, dsn string) (*sql.DB, error) {
	log.Println("open db", dsn)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, err
	}

	// 校验数据库连接
	err = db.Ping()

	return db, err
}
