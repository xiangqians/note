// mapper
// @author xiangqian
// @date 12:01 2023/05/07
package db

import (
	"database/sql"
	"fmt"
	"note/src/util/str"
	"reflect"
)

// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出LIFO

// rowsMapper 行映射
// 支持：
// 1、基本数据类型映射
// 2、结构体类型映射
// 3、基本数据类型/结构体 切片映射
func RowsMapper[T any](rows *sql.Rows) (T, int64, error) {
	var t T

	// rows is nil ?
	if rows == nil {
		return t, 0, nil
	}

	// defer close rows
	defer rows.Close()

	// 获取数据库字段名称集
	cols, err := rows.Columns()
	if err != nil {
		return t, 0, err
	}

	var count int64
	rflVal := reflect.ValueOf(&t).Elem()
	rflTyp := rflVal.Type()
	switch rflTyp.Kind() {
	// 结构体
	case reflect.Struct:
		if rows.Next() {
			count++
			dest := getDest(cols, rflTyp, rflVal)
			err = rows.Scan(dest...)
		}

	// 切片
	case reflect.Slice:
		// 创建切片
		// len 0, cap ?
		i := reflect.MakeSlice(rflTyp, 1, 1).Interface()
		t, _ = i.(T)
		rflVal = reflect.ValueOf(&t).Elem()

		// 切片长度（Len）
		l := rflVal.Len()
		// 获取切片元素
		eRflVal := rflVal.Index(0) // 获取切片第一个元素值（Value）
		eRflType := eRflVal.Type() // 切片元素类型
		e := eRflVal.Interface()   // 切片元素值
		// 切片元素类型
		switch eRflVal.Kind() {
		// Map
		case reflect.Map:
			// 尚未实现
			panic("暂不支持map切片")

		// 结构体切片、基本类型切片
		default:
			// 元素反射值
			idx := 0
			for rows.Next() {
				count++
				if idx < l {
					eRflVal = rflVal.Index(idx).Addr().Elem()
				} else {
					// reflect.New(): 返回指定类型的指针，该指针指向新创建的对象，返回指定类型指针的反射对象（Value结构体）
					e = reflect.New(eRflType).Interface()
					eRflVal = reflect.ValueOf(e).Elem()
				}
				dest := getDest(cols, eRflType, eRflVal)
				err = rows.Scan(dest...)
				if err != nil {
					return t, count, err
				}

				// 切片（slice）扩容
				if idx >= l {
					oldRflVal := rflVal
					rflVal = reflect.Append(rflVal, eRflVal)
					oldRflVal.Set(rflVal)
					rflVal = oldRflVal
				}
				idx++
			}
		}

	// Map
	case reflect.Map:
		//reflect.MakeMap(): 创建map
		panic("暂不支持map数据")

	// 普通指针类型
	default:
		if rows.Next() {
			count++
			err = rows.Scan(&t)
		}
	}

	return t, count, err
}

func getDest(cols []string, rflType reflect.Type, rflVal reflect.Value) []any {
	// kind
	switch rflType.Kind() {
	// 基本数据类型
	case reflect.Int:
		fallthrough // 执行穿透
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.String:
		return getBasicDest(cols, rflType, rflVal)

	// 结构体
	case reflect.Struct:
		return getStructDest(cols, rflType, rflVal)

	// default
	default:
		panic(fmt.Sprintf("%s '%s' %s", "不支持此类型", rflType.Kind().String(), "dest"))
	}
}

func getBasicDest(cols []string, rflType reflect.Type, rflVal reflect.Value) []any {
	// dest
	dest := make([]any, len(cols))

	// addr
	dest[0] = rflVal.Addr().Interface()

	return dest
}

func getStructDest(cols []string, rflType reflect.Type, rflVal reflect.Value) []any {
	// dest
	dest := make([]any, len(cols))

	// field
	for fi, fl := 0, rflType.NumField(); fi < fl; fi++ {
		typeField := rflType.Field(fi)

		// 兼容 FieldAlign() int （如果是struct字段，对齐后占用的字节数）
		if typeField.Type.Kind() == reflect.Struct {
			for sfi, sfl := 0, typeField.Type.NumField(); sfi < sfl; sfi++ {
				setStructDest(cols, &dest, typeField.Type.Field(sfi), rflVal)
			}
		} else {
			setStructDest(cols, &dest, typeField, rflVal)
		}
	}

	return dest
}

func setStructDest(cols []string, dest *[]any, typeField reflect.StructField, rflVal reflect.Value) {
	name := typeField.Tag.Get("sql")
	if name == "" {
		name = str.ConvHumpToUnderline(typeField.Name)
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
