// @author xiangqian
// @date 20:47 2023/06/10
package dbnew

type Db interface {
	// Begin 开启事务
	Begin() (err error)

	// Add 新增
	// returns the integer generated by the database in response to a command.
	// Typically this will be from an "auto increment" column when inserting a new row.
	// Not all databases support this feature, and the syntax of such statements varies.
	// return insertId
	Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error)

	// Del 删除
	// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
	// return affect
	Del(sql string, args ...any) (rowsAffected int64, err error)

	// Upd 更新
	// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
	Upd(sql string, args ...any) (rowsAffected int64, err error)

	// Get 查询
	Get(sql string, args ...any) (result Result, err error)

	// Page 分页查询
	Page(sql string, current int64, size uint8, args ...any) (result Result, err error)

	// Commit 提交事务
	Commit() (err error)

	// Rollback 回滚事务
	Rollback() (err error)

	// Close 关闭资源
	Close() (err error)
}

// Result 查询结果
type Result interface {
	// Count 计数
	Count() (count int64)

	// Scan 扫描数据
	Scan(dest any) (err error)
}

// DbConnPool 数据库连接池
type DbConnPool interface {
	// Get 获取数据库连接
	Get() (db Db, err error)
}
