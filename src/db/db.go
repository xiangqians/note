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
	"fmt"
	"log"
	"note/src/model"
	util_crypto_md5 "note/src/util/crypto/md5"
	"reflect"
	"sort"
	"strings"
	"time"

	// https://pkg.go.dev/gorm.io/gorm
	// https://gorm.io/gen/index.html
	"gorm.io/gorm"

	// Sqlite driver based on CGO
	"gorm.io/driver/sqlite"
)

// db map
var dbMap map[string]*gorm.DB

func init() {
	// len 0, cap ?
	// cap ?
	dbMap = make(map[string]*gorm.DB, 16)
}

// Db Database
func Db(dsn string) (*gorm.DB, error) {
	key := util_crypto_md5.Encrypt([]byte(dsn), nil)
	if db, ok := dbMap[key]; ok {
		return db, nil
	}

	db, err := open(dsn)
	if err != nil {
		return nil, err
	}

	dbMap[key] = db
	return db, nil
}

// 打开数据库连接
// https://gorm.io/zh_CN/docs
// https://gorm.io/zh_CN/docs/connecting_to_the_database.html#SQLite
func open(dsn string) (db *gorm.DB, err error) {
	log.Println("open db", dsn)

	// open
	dialector := sqlite.Open(dsn)
	db, err = gorm.Open(dialector, &gorm.Config{
		// 全局模式：执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		PrepareStmt: true,
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

// Exec 执行插入、删除等
// 使用gorm框架执行原生sql：（两种方式）
// 1、gorm.DB.Exec("sql语句") 执行插入、删除等操作
// 2、gorm.DB.Raw("sql语句")  执行查询
// gorm中exec和raw方法的区别是：Raw用来查询，执行其他操作用Exec。
// (*gorm.DB).Exec does not return an error, if you want to see if your query failed or not read up on error handling with gorm. Use Exec when you don’t care about output, use Raw when you do care about the output.
func Exec(db *gorm.DB, sql string, values ...any) (rowsAffected int64, err error) {
	db = db.Exec(sql, values...)
	rowsAffected = db.RowsAffected
	err = db.Error
	return
}

// Raw 执行查询
func Raw[T any](db *gorm.DB, sql string, values ...any) (T, error) {
	var data T
	rflTyp := reflect.ValueOf(&data).Elem().Type()
	switch rflTyp.Kind() {
	// int类型
	case reflect.Int, reflect.Int8, reflect.Int16 | reflect.Int32 | reflect.Int64 |
		reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64:

	// float类型
	case reflect.Float32, reflect.Float64:

	// string类型
	case reflect.String:

	// 结构体类型
	case reflect.Struct:

	// 切片类型
	case reflect.Slice:

	default:
		panic("不支持此类型查询：" + rflTyp.Name())
	}

	// Scan? Take?
	db = db.Raw(sql, values...).Scan(&data)
	return data, db.Error
}

// Page 分页查询
// db 数据库
// sql SQL语句
// current 当前页
// size 页数量
func Page[T any](db *gorm.DB, current int64, size uint8, sql string, values ...any) (model.Page[T], error) {
	page := model.Page[T]{
		Current: current,
		Size:    size,
	}

	// 总数
	var total int64
	index := strings.Index(sql, "FROM")
	db.Raw(fmt.Sprintf("SELECT COUNT(1) %s", sql[index:]), values...).Count(&total)
	page.Total = total
	if total == 0 {
		return page, nil
	}

	// 总页数
	pageCount := total / int64(size)
	if total%int64(size) != 0 {
		pageCount += 1
	}
	page.PageCount = pageCount

	// 页数索引集
	if current == 1 || current > pageCount {
		pageIndexes := make([]int64, 0, 8)
		var pageIndex int64 = 1
		count := cap(pageIndexes)
		for {
			count--
			if count < 0 || pageIndex > pageCount {
				break
			}
			pageIndexes = append(pageIndexes, pageIndex)
			pageIndex++
		}

		length := len(pageIndexes)
		if pageIndexes[length-1] != pageCount {
			pageIndexes[length-2] = 0
			pageIndexes[length-1] = pageCount
		}
		page.PageIndexes = pageIndexes

	} else if current == pageCount {
		pageIndexes := make([]int64, 0, 8)
		var pageIndex int64 = pageCount
		count := cap(pageIndexes)
		for {
			count--
			if count < 0 || pageIndex <= 0 {
				break
			}
			pageIndexes = append(pageIndexes, pageIndex)
			pageIndex--
		}

		// 排序：升序
		sort.Slice(pageIndexes, func(i, j int) bool {
			return i > j
		})

		if pageIndexes[0] != 1 {
			pageIndexes[0] = 1
			pageIndexes[1] = 0
		}
		page.PageIndexes = pageIndexes

	} else {
		pageIndexes := make([]int64, 0, 6+1+6)
		var pageIndex int64 = current - 6
		if pageIndex <= 0 {
			pageIndex = 1
		}
		index := 0 // 当前页索引在数组中位置
		count := cap(pageIndexes)
		for {
			count--
			if count < 0 || pageIndex > pageCount {
				break
			}
			pageIndexes = append(pageIndexes, pageIndex)
			if current == pageIndex {
				index = len(pageIndexes) - 1
			}
			pageIndex++
		}

		length := len(pageIndexes)
		// ... 在右侧
		if pageIndexes[0] == 1 && index < 4 {
			if length >= 8 && pageIndexes[8-1] != pageCount {
				pageIndexes[8-2] = 0
				pageIndexes[8-1] = pageCount
				pageIndexes = pageIndexes[0:8]
			}

		} else
		// ... 在左侧
		if length >= 8 && pageIndexes[length-1] == pageCount && index >= length-4 {
			if pageIndexes[length-8] != 1 {
				pageIndexes[length-8] = 1
				pageIndexes[length-8+1] = 0
				pageIndexes = pageIndexes[length-8:]
			}
		} else
		// ... 在左右两侧
		if length > 8 {
			pageIndexes = pageIndexes[index-4 : index+4+1]
			length = len(pageIndexes)
			pageIndexes[0] = 1
			pageIndexes[1] = 0
			pageIndexes[length-2] = 0
			pageIndexes[length-1] = pageCount
		}
		page.PageIndexes = pageIndexes
	}

	// 数据
	var data []T
	rflTyp := reflect.ValueOf(&data).Elem().Type()
	// 创建切片：len 0, cap ?
	i := reflect.MakeSlice(rflTyp, 0, int(size)).Interface()
	data = i.([]T)
	offset := (current - 1) * int64(size)
	limit := size
	db = db.Raw(fmt.Sprintf("%s LIMIT %d,%d", sql, offset, limit), values...)
	err := db.Scan(&data).Error
	page.Data = data
	return page, err
}
