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
// @author xiangqian
// @date 20:47 2023/06/10
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"note/src/model"
	"reflect"
	"strings"
)

var db *sql.DB

// 初始化数据库&连接池
func init() {
	// 驱动名（driver name）
	driver := model.Ini.Db.Driver

	// 数据源（data source name）
	dsn := model.Ini.Db.Dns

	var err error
	db, err = sql.Open(driver, dsn)
	log.Printf("open %s db: %s\n", driver, dsn)
	if err != nil {
		panic(err)
	}

	//type sql.DB struct {
	//	// Atomic access only. At top of struct to prevent mis-alignment
	//	// on 32-bit platforms. Of type time.Duration.
	//	waitDuration int64 // Total time waited for new connections. 等待新连接的总时间，用于统计
	//
	//	connector driver.Connector // 由数据库驱动实现的连接器
	//	// numClosed is an atomic counter which represents a total number of
	//	// closed connections. Stmt.openStmt checks it before cleaning closed
	//	// connections in Stmt.css.
	//	numClosed uint64 // 关闭连接数
	//
	//	mu           sync.Mutex    // protects following fields.
	//	freeConn     []*driverConn // free connections ordered by returnedAt oldest to newest. 可用连接（池）
	//	connRequests map[uint64]chan connRequest // 连接请求表，key是分配的自增键
	//	nextRequest  uint64 // Next key to use in connRequests.  连接请求自增键
	//	numOpen      int    // number of opened and pending open connections. 已打开 + 即将打开的连接数
	//	// Used to signal the need for new connections
	//	// a goroutine running connectionOpener() reads on this chan and
	//	// maybeOpenNewConnections sends on the chan (one send per needed connection)
	//	// It is closed during db.Close(). The close tells the connectionOpener
	//	// goroutine to exit.
	//	openerCh          chan struct{} // 需要新的连接
	//	closed            bool
	//	dep               map[finalCloser]depSet
	//	lastPut           map[*driverConn]string // stacktrace of last conn's put; debug only
	//	maxIdleCount      int                    // zero means defaultMaxIdleConns; negative means 0
	//	maxOpen           int                    // <= 0 means unlimited
	//	maxLifetime       time.Duration          // maximum amount of time a connection may be reused
	//	maxIdleTime       time.Duration          // maximum amount of time a connection may be idle before being closed
	//	cleanerCh         chan struct{}
	//	waitCount         int64 // Total number of connections waited for.
	//	maxIdleClosed     int64 // Total number of connections closed due to idle count.
	//	maxIdleTimeClosed int64 // Total number of connections closed due to idle time.
	//	maxLifetimeClosed int64 // Total number of connections closed due to max connection lifetime limit.
	//
	//	stop func() // stop cancels the connection opener.
	//}

	// driverConn wraps a driver.Conn with a mutex, to be held during all calls into the Conn.
	// (including any calls onto interfaces returned via that Conn, such as calls on Tx, Stmt, Result, Rows)
	//type sql.driverConn struct {
	//	db        *DB // 数据库句柄
	//	createdAt time.Time
	//
	//	sync.Mutex  // guards following. 锁
	//	ci          driver.Conn // 对应具体的连接
	//	needReset   bool // The connection session should be reset before use if true.
	//	closed      bool // 是否已关闭
	//	finalClosed bool // ci.Close has been called. 是否最终关闭
	//	openStmt    map[*driverStmt]bool // 在这个连接上打开的状态
	//
	//	// guarded by db.mu
	//	inUse      bool // 连接是否占用
	//	returnedAt time.Time // Time the connection was created or returned.
	//	onPut      []func()  // code (with db.mu held) run when conn is next returned. 连接归还时要运行的函数，在 noteUnusedDeriverStatement 添加
	//	dbmuClosed bool      // same as closed, but guarded by db.mu, for removeClosedStmtLocked. 和 closed 状态一致，但是由锁保护，用于 removeClosedStmtLocked
	//}

	// 在Golang中，database/sql包已经集成了连接池的功能

	// sql.DB连接池是如何工作的呢？
	// sql.DB连接池包含两种类型的连接：”正在使用“连接和“空闲”连接。
	// 当使用连接执行数据库任务（例如执行SQL语句或查询行）时，该连接被标记为”正在使用“，任务完成后，该连接被标记为“空闲”。
	//
	// 当使用执行数据库操作（例如Exec，Query）时，首先检查池中是否有可用的“空闲”连接。
	// 如果有可用的连接，那么将重用这个现有连接，并在任务期间将其标记为”正在使用“。
	// 如果在您需要“空闲”连接时池中没有“空闲”连接，那么将创建一个新的连接。

	// 连接失败
	// 不必检查或者尝试处理连接失败的情况。
	// 当进行数据库操作的时候，如果连接失败了，database/sql会处理。
	// 实际上，当从连接池取出的连接断开的时候，database/sql会自动尝试重连10次。
	// 仍然无法重连的情况下会自动从连接池再获取一个或者新建另外一个。

	// 设置池中“打开”连接（”正在使用“连接和“空闲”连接）数量的上限。默认情况下，打开的连接数是无限的。
	// 注意 “打开”连接 包括 ”正在使用“连接和“空闲”连接，不仅仅是“正在使用”连接。
	// 一般来说，MaxOpenConns设置得越大，可以并发执行的数据库查询就越多，连接池本身成为应用程序中的瓶颈的风险就越低。
	// 但让它无限并不是最好的选择。默认情况下，PostgreSQL最多100个打开连接的硬限制，如果达到这个限制的话，它将导致pq驱动返回"sorry, too many clients already"错误。
	// 注意：最大打开连接数限制可以在postgresql.conf文件中对max_connections设置来更改。
	// 为了避免这个错误，将池中打开的连接数量限制在100以下是有意义的，可以为其他需要使用PostgreSQL的应用程序或会话留下足够的空间。
	// 设置MaxOpenConns限制的另一个好处是，它充当一个非常基本的限流器，防止数据库同时被大量任务压垮。
	// 如果达到MaxOpenConns限制，并且所有连接都在使用中，那么任何新的数据库任务将被迫等待，直到有连接空闲。
	// 在我们的API上下文中，用户的HTTP请求可能在等待空闲连接时无限期地“挂起”。因此，为了缓解这种情况，使用上下文为数据库任务设置超时是很重要的。
	// sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(model.Ini.Db.MaxOpenConns) // 应该根据基准测试和压测结果调整这个值

	// 设置一个连接保持可用的最长时间。默认连接的存活时间没有限制，永久可用。
	// 如果设置ConnMaxLifetime的值为1小时，意味着所有的连接在创建后，经过一个小时就会被标记为失效连接，标志后就不可复用。
	// 但需要注意：
	// 1、这并不能保证一个连接将在池中存在一整个小时；有可能某个连接由于某种原因变得不可用，并在此之前自动关闭。
	// 2、连接在创建后一个多小时内仍然可以被使用—只是在这个时间之后它不能被重用。
	// 3、这不是一个空闲超时。连接将在创建后一小时过期，而不是在空闲后一小时过期。
	// 4、Go每秒运行一次后台清理操作，从池中删除过期的连接。
	// 理论上，ConnMaxLifetime为无限大（或设置为很长生命周期）将提升性能，因为这样可以减少新建连接。但是在某些情况下，设置短期存活时间有用。比如：
	// 1、如果SQL数据库对连接强制设置最大存活时间，这时将ConnMaxLifetime设置稍短时间更合理。
	// 2、有助于数据库替换（优雅地交换数据库）
	// 如果决定对连接池设置ConnMaxLifetime，那么一定要记住连接过期(然后重新创建)的频率。例如，如果连接池中有100个打开的连接，而ConnMaxLifetime为1分钟，那么应用程序平均每秒可以杀死并重新创建多达1.67个连接。频率太大而最终影响性能。
	// sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(model.Ini.Db.ConnMaxLifetime)

	// 设置池中“空闲”连接数的上限。缺省情况下，最大空闲连接数为2。
	// 理论上，在池中允许更多的空闲连接将增加性能。因为它减少了从头建立新连接发生概率，因此有助于节省资源。
	// 但要意识到保持空闲连接是有代价的。它占用了本来可以用于应用程序和数据库的内存，而且如果一个连接空闲时间过长，它也可能变得不可用。例如，默认情况下MySQL会自动关闭任何8小时未使用的连接。
	// 因此，与使用更小的空闲连接池相比，将MaxIdleConns设置得过高可能会导致更多的连接变得不可用，浪费资源。因此保持适量的空闲连接是必要的。理想情况下，你只希望保持一个连接空闲，可以快速使用。
	// 另一件要指出的事情是MaxIdleConns值应该总是小于或等于MaxOpenConns。Go会强制保证这点，并在必要时自动减少MaxIdleConns值。
	// sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(model.Ini.Db.MaxIdleConns)

	// SetConnMaxIdleTime()方法在Go 1.15版本引入对ConnMaxIdleTime进行配置。
	// 其效果和ConnMaxLifeTime类似，但这里设置的是：在被标记为失效之前一个连接最长空闲时间。
	// 例如，如果我们将ConnMaxIdleTime设置为1小时，那么自上次使用以后在池中空闲了1小时的任何连接都将被标记为过期并被后台清理操作删除。
	// 这个配置非常有用，因为它意味着我们可以对池中空闲连接的数量设置相对较高的限制，但可以通过删除不再真正使用的空闲连接来周期性地释放资源。
	db.SetConnMaxIdleTime(model.Ini.Db.ConnMaxIdleTime)

	// 因为每一个连接都是惰性创建的，如何验证sql.Open调用之后，sql.DB对象可用呢？
	// 通常使用sql.Ping()方法初始化，调用完毕后会马上把连接返回给连接池。
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

// Get 获取数据库连接
func Get() *Db {
	return &Db{db: db}
}

// Stats DB状态
func Stats() string {
	stats := db.Stats()
	statsString := fmt.Sprint("\n\t最大连接数\t\t:", stats.MaxOpenConnections,
		"\n\t连接池状态",
		"\n\t\t当前连接数\t:", stats.OpenConnections, "（”正在使用“连接和“空闲”连接）",
		"\n\t\t正在使用连接数\t:", stats.InUse,
		"\n\t\t空闲连接数\t:", stats.Idle,
		"\n\t统计",
		"\n\t\t等待连接数\t:", stats.WaitCount,
		"\n\t\t等待创建新连接时长（秒）:", stats.WaitDuration.Seconds(),
		"\n\t\t空闲超限关闭数\t:", stats.MaxIdleClosed, "（达到MaxIdleConns而关闭的连接数量）",
		"\n\t\t空闲超时关闭数\t:", stats.MaxIdleTimeClosed,
		"\n\t\t连接超时关闭数\t:", stats.MaxLifetimeClosed, "（达到ConnMaxLifetime而关闭的连接数量）")
	return statsString
}

type Db struct {
	db *sql.DB // db
	tx *sql.Tx // 事务支持
}

// Begin 开启事务
// sql.Begin()调用完毕后将连接传递给sql.Tx类型对象，当Commit()或Rollback()方法调用后释放连接
func (db *Db) Begin() (err error) {
	db.tx, err = db.db.Begin()
	return
}

// Add 新增
// returns the integer generated by the database in response to a command.
// Typically this will be from an "auto increment" column when inserting a new row.
// Not all databases support this feature, and the syntax of such statements varies.
// return insertId
func (db *Db) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	log.Println("[?ms]", sql, Stats())
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	insertId, _ = result.LastInsertId()
	return
}

// Del 删除
// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
// return affect
func (db *Db) Del(sql string, args ...any) (rowsAffected int64, err error) {
	log.Println("[?ms]", sql, Stats())
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return
}

// Upd 更新
// returns the number of rows affected by an update, insert, or delete. Not every database or database driver may support this.
func (db *Db) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	log.Println("[?ms]", sql, Stats())
	result, err := db.exec(sql, args...)
	if err != nil {
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return
}

// sql.Exec()调用完毕后会马上把连接返回给连接池，但是它返回的Result对象还保留这连接的引用，当后面的代码需要处理结果集的时候连接将会被重用
func (db *Db) exec(sql string, args ...any) (sql.Result, error) {
	if db.tx != nil {
		return db.tx.Exec(sql, args...)
	}
	return db.db.Exec(sql, args...)
}

// Get 查询
func (db *Db) Get(sql string, args ...any) (*Result, error) {
	log.Println("[?ms]", sql, Stats())
	rows, err := db.query(sql, args...)
	return &Result{rows: rows}, err
}

// sql.Query()调用完毕后会将连接传递给sql.Rows类型，当然后者迭代完毕或者显示的调用Close()方法后，连接将会被释放回到连接池
// sql.QueryRow()调用完毕后会将连接传递给sql.Row类型，当.Scan()方法调用之后把连接释放回到连接池
func (db *Db) query(sql string, args ...any) (*sql.Rows, error) {
	if db.tx != nil {
		return db.tx.Query(sql, args...)
	}
	return db.db.Query(sql, args...)
}

// Page 分页查询
func (db *Db) Page(sql string, current int64, size uint8, args ...any) (*Result, error) {
	// 计数
	result, err := db.Get(fmt.Sprintf("SELECT COUNT(1) FROM (%s) r", sql), args...)
	if err != nil {
		return &Result{}, err
	}
	var count int64
	result.Scan(&count)

	// 查询分页数据
	offset := (current - 1) * int64(size)
	limit := size
	result, err = db.Get(fmt.Sprintf("%s LIMIT %d,%d", sql, offset, limit), args...)

	result.count = count
	return result, err
}

// Commit 提交事务
func (db *Db) Commit() error {
	if db.tx != nil {
		return db.tx.Commit()
	}
	return nil
}

// Rollback 回滚事务
func (db *Db) Rollback() error {
	if db.tx != nil {
		return db.tx.Rollback()
	}
	return nil
}

// Close 关闭资源
func (db *Db) Close() error {
	if db.db != nil {
		err := db.db.Close()
		db.db = nil
		return err
	}
	return nil
}

// Result 查询结果
type Result struct {
	count int64
	rows  *sql.Rows
}

// Count 计数
func (result *Result) Count() int64 {
	return result.count
}

// Scan 扫描数据
func (result *Result) Scan(dest any) error {
	if result.rows == nil {
		return nil
	}

	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO
	defer result.rows.Close()

	v := reflect.ValueOf(dest).Elem()
	switch v.Kind() {
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
		if result.next() {
			return result.scanStruct(dest)
		}

	// 切片
	case reflect.Slice:
		t := v.Type().Elem()
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

func (result *Result) next() bool {
	return result.rows.Next()
}

// 扫描基本数据类型
func (result *Result) scanDefault(dest any) error {
	return result.rows.Scan(dest)
}

// 扫描结构体
func (result *Result) scanStruct(dest any) error {
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

	v := reflect.ValueOf(dest).Elem()
	t := v.Type()
	result.dest(&cols, &newDest, t, v)

	return result.rows.Scan(newDest...)
}

func (result *Result) dest(cols *[]string, dest *[]any, t reflect.Type, v reflect.Value) {
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
