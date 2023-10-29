// @author xiangqian
// @date 19:10 2023/10/29
package dbnew

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	util_crypto_md5 "note/src/util/crypto/md5"
	"strings"
	"sync"
	"time"
)

type GormDb struct {
	db *gorm.DB
}

func (db *GormDb) Begin() (err error) {
	return
}

// Add
// Exec 执行插入、删除等
// 使用gorm框架执行原生sql：（两种方式）
// 1、gorm.DB.Exec("sql语句") 执行插入、删除等操作
// 2、gorm.DB.Raw("sql语句")  执行查询
// gorm中exec和raw方法的区别是：Raw用来查询，执行其他操作用Exec。
// (*gorm.DB).Exec does not return an error, if you want to see if your query failed or not read up on error handling with gorm. Use Exec when you don’t care about output, use Raw when you do care about the output.
func (db *GormDb) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	// gorm Exec不支持获取自增id
	// 但是gorm orm是对底层database/sql的封装，所以进行降级执行

	sqlDb, err := db.db.DB()
	if err != nil {
		return
	}

	result, err := sqlDb.Exec(sql, args...)
	if err != nil {
		return
	}

	rowsAffected, _ = result.RowsAffected()
	insertId, _ = result.LastInsertId()
	return
}

func (db *GormDb) Del(sql string, args ...any) (rowsAffected int64, err error) {
	tx := db.db.Exec(sql, args...)
	rowsAffected = tx.RowsAffected
	err = tx.Error
	return
}

func (db *GormDb) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	tx := db.db.Exec(sql, args...)
	rowsAffected = tx.RowsAffected
	err = tx.Error
	return
}

func (db *GormDb) Get(sql string, args ...any) (Result, error) {
	tx := db.db.Raw(sql, args...)
	return GormResult{tx: tx}, tx.Error
}

func (db *GormDb) Page(sql string, current int64, size uint8, args ...any) (Result, error) {
	// 计数
	index := strings.Index(sql, "FROM")
	result, err := db.Get(fmt.Sprintf("SELECT COUNT(1) %s", sql[index:]), args...)
	if err != nil {
		return GormResult{}, err
	}
	var count int64
	result.Scan(&count)

	// 查询分页数据
	offset := (current - 1) * int64(size)
	limit := size
	tx := db.db.Raw(fmt.Sprintf("%s LIMIT %d,%d", sql, offset, limit), args...)
	return GormResult{tx: tx, count: count}, tx.Error
}

func (db *GormDb) Commit() (err error) {
	return
}

func (db *GormDb) Rollback() (err error) {
	return
}

func (db *GormDb) Close() (err error) {
	return
}

type GormResult struct {
	tx    *gorm.DB
	count int64
}

func (result GormResult) Count() int64 {
	return result.count
}

func (result GormResult) Scan(dest any) error {
	if result.tx == nil {
		return nil
	}

	// Scan? Take?
	return result.tx.Scan(dest).Error
}

type GormDbConnPool struct {
	// sync.Mutex 是一个基本的同步原语，可以实现并发环境下的线程安全
	mutex sync.Mutex

	// db map
	m map[string]*gorm.DB
}

func (dbConnPool *GormDbConnPool) Get(dsn string) (Db, error) {
	dbConnPool.mutex.Lock() // 获取锁
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer dbConnPool.mutex.Unlock() // 释放锁

	if dbConnPool.m == nil {
		// len 0, cap ?
		// cap ?
		dbConnPool.m = make(map[string]*gorm.DB, 16)
	}

	key := util_crypto_md5.Encrypt([]byte(dsn), nil)
	if db, ok := dbConnPool.m[key]; ok {
		return &GormDb{db: db}, nil
	}

	db, err := dbConnPool.open(dsn)
	if err != nil {
		return nil, err
	}

	dbConnPool.m[key] = db
	return &GormDb{db: db}, nil
}

// 打开数据库连接
// https://gorm.io/zh_CN/docs
// https://gorm.io/zh_CN/docs/connecting_to_the_database.html#SQLite
// dsn : DataSourceName
func (dbConnPool *GormDbConnPool) open(dsn string) (db *gorm.DB, err error) {
	log.Println("open db", dsn)

	dialector := sqlite.Open(dsn)
	db, err = gorm.Open(dialector, &gorm.Config{
		// 全局模式：执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		PrepareStmt: true,

		// 日志记录器
		// Gorm有一个默认logger实现，默认情况下，它会打印慢SQL（默认200ms）和错误
		Logger: logger.New(log.New(log.Writer(), "", log.LstdFlags|log.LstdFlags), logger.Config{
			// 设定慢查询时间阈值为1ns
			//SlowThreshold: 1 * time.Nanosecond,
			// 设置日志级别，只有Info和Warn级别会输出慢查询日志
			LogLevel: logger.Info,
			// 忽略找不到记录错误
			IgnoreRecordNotFoundError: false,
			// 彩色日志输出
			Colorful: false,
		}),
	})
	if err != nil {
		return
	}

	// 配置连接池
	// 通过数据库连接池，我们可以避免频繁创建和销数据库连接所带来的开销，GROM的数据连接池底层是通过database/sql来实现的，所以其设置方法与database/sql是一样的。
	sqlDb, err := db.DB()
	if err != nil {
		return
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Minute * 10)
	return
}
