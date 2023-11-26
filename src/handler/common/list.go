// @author xiangqian
// @date 23:21 2023/10/23
package common

import (
	"fmt"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func List(request *http.Request, session *session.Session, table string) (string, model.Response) {
	// html模板
	html := func(page model.Page, err any) (string, model.Response) {
		return fmt.Sprintf("%s/list", table), model.Response{
			Msg:  util_string.String(err),
			Data: model.Page{},
		}
	}

	// 当前页
	current, _ := strconv.ParseInt(request.URL.Query().Get("current"), 10, 64)
	if current <= 0 {
		current = 1
	}

	// 页数量
	size64, _ := strconv.ParseInt(request.URL.Query().Get("size"), 10, 64)
	size := uint8(size64)
	if size <= 0 || size > 100 {
		size = 10
	}

	// 检索条件
	search := strings.TrimSpace(request.URL.Query().Get("search"))
	search = "name:12342d type:2df23 size>=1234"

	page := model.Page{
		Current: current,
		Size:    size,
		Total:   0,
		Search:  search,
	}

	search = fmt.Sprintf(" %s", search)

	columns := []column{
		{
			name:      " id:",
			index:     -1,
			kind:      reflect.Int64,
			statement: "`id` = ?",
			value:     nil,
		},
		{
			name:      " name:",
			index:     -1,
			kind:      reflect.String,
			statement: "`name` LIKE '%' || ? || '%'", // sqlite在模糊查询时大小写不敏感
			value:     nil,
		},
		{
			name:      " type:",
			index:     -1,
			kind:      reflect.String,
			statement: "`type` = ?",
			value:     nil,
		},
		{
			name:      " size>=",
			index:     -1,
			kind:      reflect.Int64,
			statement: "`size` >= ?",
			value:     nil,
		},
		{
			name:      " del:",
			index:     -1,
			kind:      reflect.Uint8,
			statement: "`del` = ?",
			value:     0,
		},
	}

	length := len(columns)

	if search != "" {
		// 获取字段所在检索文本search中的位置
		for i := 0; i < length; i++ {
			index := strings.Index(search, columns[i].name)
			(&columns[i]).index = index
		}

		// 字段切片根据索引（index）升序排序
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].index < columns[j].index
		})

		// 获取字段值
		for i := 0; i < length; i++ {
			index := columns[i].index
			if index == -1 {
				continue
			}

			var nextIndex int = -1
			for i := i + 1; i < length; i++ {
				nextIndex = columns[i].index
				if nextIndex != -1 {
					break
				}
			}

			if nextIndex == -1 {
				nextIndex = len(search)
			}

			value := strings.TrimSpace(search[index+len(columns[i].name) : nextIndex])
			switch columns[i].kind {
			case reflect.Uint8:
				i, _ := strconv.ParseInt(value, 10, 64)
				(&columns[i]).value = uint8(i)

			case reflect.Int64:
				(&columns[i]).value, _ = strconv.ParseInt(value, 10, 64)

			case reflect.String:
				(&columns[i]).value = value
			}
		}

	}

	// len 0, cap ?
	statements := make([]string, 0, length)
	values := make([]any, 0, length)

	for i := 0; i < length; i++ {
		if columns[i].value != nil {
			statements = append(statements, columns[i].statement)
			values = append(values, columns[i].value)
		}
	}

	db := db.Get()
	sql := fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `history`, `history_size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE %s",
		table,
		strings.Join(statements, " AND "))
	result, err := db.Page(sql, current, size, values...)
	if err != nil {
		return html(page, err)
	}

	page.Total = result.Count()
	(&page).InitIndexes()

	var data any
	switch table {
	case "image":
		data = []model.Image{}
	}
	err = result.Scan(&data)
	page.Data = data
	return html(page, err)
}

// 字段
type column struct {
	name      string       // 字段名集
	index     int          // 字段在检索字符串（search）中索引
	kind      reflect.Kind // 字段类型
	statement string       // 语句
	value     any          // 值
}

//// RedirectToList 重定向到列表
//// ctx  : *gin.Context
//// table: 数据表名
//// msg  : 重定向消息（没有消息就是最好的消息）
//func RedirectToList[T any](ctx *gin.Context, msg any) {
//	// 数据表名
//	table := Table[T]()
//
//	// 获取分页参数
//	current, _ := context.Query[int64](ctx, "current")
//	size, _ := context.Query[uint8](ctx, "size")
//	search, _ := context.Query[string](ctx, "search")
//
//	// 重定向到图片首页
//	context.Redirect(ctx, fmt.Sprintf("/%s", table), map[string]any{
//		"current": current,
//		"size":    size,
//		"search":  search,
//	}, msg)
//}
