// DB
// https://pkg.go.dev/github.com/mattn/go-sqlite3
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"note/src/typ"
	"note/src/util"
	"reflect"
	"strings"
)

func Page[T any](dsn string, pageReq typ.PageReq, sql string, args ...any) (typ.Page[T], error) {
	current := pageReq.Current
	size := pageReq.Size
	page := typ.Page[T]{
		Current: current,
		Size:    size,
	}

	// 总记录数
	_sql := fmt.Sprintf("SELECT COUNT(1) %s", sql[strings.Index(sql, "FROM"):])
	if strings.Contains(_sql, "GROUP BY") {
		_sql = fmt.Sprintf("SELECT COUNT(1) FROM (%s) r", _sql)
	}
	total, _, err := Qry[int64](dsn, _sql, args...)
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
	data, count, err := Qry[[]T](dsn, sql, args...)
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
	_db, err := db(dsn)
	if err != nil {
		return t, 0, err
	}
	defer _db.Close()

	// query
	rows, err := _db.Query(sql, args...)
	if err != nil {
		return t, 0, err
	}
	defer rows.Close()

	// 通过反射初始化实例
	rflVal := reflect.ValueOf(t)
	switch rflVal.Type().Kind() {
	case reflect.Slice:
		i := reflect.MakeSlice(rflVal.Type(), 1, 1).Interface()
		t, _ = i.(T)
	}

	// 行映射
	count, err := rowsMapper(rows, &t)
	if err != nil {
		return t, count, err
	}

	return t, count, nil
}

// Add return insertId
func Add(dsn string, sql string, args ...any) (int64, error) {
	_db, err := db(dsn)
	if err != nil {
		return 0, err
	}
	defer _db.Close()

	res, err := _db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Upd(dsn string, sql string, args ...any) (int64, error) {
	return exec(dsn, sql, args...)
}

func Del(dsn string, sql string, args ...any) (int64, error) {
	return exec(dsn, sql, args...)
}

// return affect
func exec(dsn string, sql string, args ...any) (int64, error) {
	_db, err := db(dsn)
	if err != nil {
		return 0, err
	}
	defer _db.Close()

	res, err := _db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// db
// dsn: DataSourceName
func db(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}

// 字段集映射
// 支持 1）一个或多个属性映射；2）结构体映射；3）结构体切片映射
func rowsMapper(rows *sql.Rows, i any) (int64, error) {
	cols, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	var record int64
	rflType := reflect.TypeOf(i).Elem()
	rflVal := reflect.ValueOf(i).Elem()
	switch rflType.Kind() {
	// 结构体
	case reflect.Struct:
		if rows.Next() {
			record++
			dest := getDest(cols, rflType, rflVal)
			err = rows.Scan(dest...)
		}

	// 切片
	case reflect.Slice:
		eVal := rflVal.Index(0)
		l := rflVal.Len()
		switch eVal.Kind() {
		// 结构体切片
		case reflect.Struct:
			e := eVal.Interface()
			eRflType := reflect.TypeOf(e)
			var eRflVal reflect.Value
			idx := 0
			for rows.Next() {
				record++
				if idx < l {
					eRflVal = rflVal.Index(idx).Addr().Elem()
				} else {
					pE := reflect.New(eRflType).Interface()
					eRflVal = reflect.ValueOf(pE).Elem()
				}
				dest := getDest(cols, eRflType, eRflVal)
				err = rows.Scan(dest...)
				if err != nil {
					return record, err
				}

				// 切片（slice）扩容
				if idx >= l {
					rflVal0 := rflVal
					rflVal = reflect.Append(rflVal, eRflVal)
					rflVal0.Set(rflVal)
					rflVal = rflVal0
				}
				idx++
			}

		// 普通指针类型数组
		default:
			if rows.Next() {
				record++
				dest := make([]any, l)
				for ei := 0; ei < l; ei++ {
					e := rflVal.Index(ei).Interface()
					dest[ei] = e
				}
				err = rows.Scan(dest...)
			}
		}

	// 普通指针类型
	default:
		if rows.Next() {
			record++
			err = rows.Scan(i)
		}
	}

	return record, err
}

func getDest(cols []string, rflType reflect.Type, rflVal reflect.Value) []any {
	dest := make([]any, len(cols))
	for fi, fl := 0, rflType.NumField(); fi < fl; fi++ {
		typeField := rflType.Field(fi)

		// 兼容 FieldAlign() int （如果是struct字段，对齐后占用的字节数）
		if typeField.Type.Kind() == reflect.Struct {
			for sfi, sfl := 0, typeField.Type.NumField(); sfi < sfl; sfi++ {
				setDest(cols, &dest, typeField.Type.Field(sfi), rflVal)
			}
		} else {
			setDest(cols, &dest, typeField, rflVal)
		}
	}

	return dest
}

func setDest(cols []string, dest *[]any, typeField reflect.StructField, rflVal reflect.Value) {
	name := typeField.Tag.Get("sql")
	if name == "" {
		name = util.NameHumpToUnderline(typeField.Name)
	}
	for ci, col := range cols {
		if col == name {
			valField := rflVal.FieldByName(typeField.Name)
			if valField.CanAddr() {
				(*dest)[ci] = valField.Addr().Interface()
			}
			break
		}
	}
}
