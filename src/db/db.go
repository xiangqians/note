/*
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
https://github.com/mattn/go-sqlite3/issues/855
https://github.com/mattn/go-sqlite3/issues/975
require (
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
)

解决方法1：拉取其他版本
https://github.com/mattn/go-sqlite3
Latest stable version is v1.14 or later, not v2.
go get github.com/mattn/go-sqlite3@v1.14.16

解决方法2：在不同系统构建不同可执行包
*/

// db
// @author xiangqian
// @date 20:47 2023/06/10
package db

import (
	"time"

	// https://pkg.go.dev/gorm.io/gorm
	// https://gorm.io/gen/index.html
	"gorm.io/gorm"

	// Sqlite driver based on CGO
	"gorm.io/driver/sqlite"
)

// db
var db *gorm.DB

// init 初始化连接
// https://gorm.io/zh_CN/docs/connecting_to_the_database.html#SQLite
func init() {
	// open
	dialector := sqlite.Open("C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db")
	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		// 全局模式：执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	// 配置连接池
	// 通过数据库连接池，我们可以避免频繁创建和销数据库连接所带来的开销，GROM的数据连接池底层是通过database/sql来实现的，所以其设置方法与database/sql是一样的。
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Minute * 10)
}

// DB Database
// https://gorm.io/zh_CN/docs
//
// 使用gorm框架执行原生sql：（两种方式）
// 1、gorm.DB.Exec("sql语句") // 执行插入删除等操作使用
// 2、gorm.DB.Raw("sql语句") // 执行查询操作时使用
// gorm中exec和raw方法的区别是：Raw用来查询，执行其他操作用Exec。
// (*gorm.DB).Exec does not return an error, if you want to see if your query failed or not read up on error handling with gorm. Use Exec when you don’t care about output, use Raw when you do care about the output.
func DB() *gorm.DB {
	return db
}
