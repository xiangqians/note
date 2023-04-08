// reflect
// @author xiangqian
// @date 22:49 2023/03/23
package reflect

import (
	"reflect"
)

// CallField 反射执行属性
// i: 实例
// name: 属性名
func CallField[T any](i any, name string) T {
	var t T
	ref := reflect.ValueOf(i)
	field := ref.FieldByName(name)
	if !field.IsValid() {
		return t
	}

	tKind := reflect.ValueOf(t).Type().Kind()
	fKind := field.Kind()
	if fKind != tKind {
		return t
	}

	switch fKind {
	case reflect.Int, reflect.Uint,
		reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32,
		reflect.Int64, reflect.Uint64,
		reflect.String:
		return field.Interface().(T)

	default:
		return t
	}

	return t
}

// CallMethod 反射执行方法
// i: 实例
// name: 属性名
// in: 入参，如果没有参数可以传 nil 或者空切片 make([]reflect.Value, 0)
func CallMethod(i any, name string, in []reflect.Value) interface{} {
	ref := reflect.ValueOf(i)
	method := ref.MethodByName(name)
	if method.IsValid() {
		r := method.Call(in)
		return r[0].Interface()
	}
	return nil
}
