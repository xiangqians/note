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

func List(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// html模板
	html := func(page model.Page, err any) (string, model.Response) {
		return fmt.Sprintf("%s/list", table),
			model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
				"table": table,
				"page":  page,
			}}
	}

	// 当前页
	current, _ := strconv.ParseInt(request.URL.Query().Get("current"), 10, 64)
	if current <= 0 {
		current = 1
	}

	// 页数量
	size64, _ := strconv.ParseInt(request.URL.Query().Get("size"), 10, 64)
	size := uint8(size64)
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 检索条件
	search := strings.TrimSpace(request.URL.Query().Get("search"))

	page := model.Page{
		Current: current,
		Size:    size,
		Search:  search,
		Total:   0,
	}

	columns := []column{
		{
			name:      " id:",
			index:     -1,
			kind:      reflect.Int64,
			statement: "`id` = ?",
			value:     nil,
		},
		{
			name:      " pid:",
			index:     -1,
			kind:      reflect.Int64,
			statement: "`pid` = ?",
			value:     int64(0),
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
			name:      " del:",
			index:     -1,
			kind:      reflect.Uint8,
			statement: "`del` = ?",
			value:     uint8(0),
		},
	}

	length := len(columns)

	if search != "" {
		search = fmt.Sprintf(" %s", search)

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
				i64, _ := strconv.ParseInt(value, 10, 64)
				(&columns[i]).value = uint8(i64)

			case reflect.Int64:
				(&columns[i]).value, _ = strconv.ParseInt(value, 10, 64)

			case reflect.String:
				(&columns[i]).value = value
			}
		}
	}

	for _, column := range columns {
		if column.name == " del:" {
			del := column.value.(uint8)
			if del != 0 && del != 1 {
				return html(page, nil)
			}
			break
		}
	}

	// len 0, cap ?
	statements := make([]string, 0, length)
	values := make([]any, 0, length)

	for i := 0; i < length; i++ {
		column := columns[i]
		if column.value != nil {
			if column.name == " pid:" && table != "note" {
				continue
			}
			statements = append(statements, column.statement)
			values = append(values, column.value)
		}
	}

	db := db.Get()
	sql := "SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time`"
	if table == "note" {
		sql += ", pid"
	}
	sql += fmt.Sprintf(" FROM `%s` ", table)
	sql += fmt.Sprintf(" WHERE %s ", strings.Join(statements, " AND "))
	switch table {
	case "image", "audio", "video":
		sql += " ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC"
	case "note":
		sql += " ORDER BY `type`, `name`, (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC"
	}
	result, err := db.Page(sql, current, size, values...)
	if err != nil {
		return html(page, err)
	}

	page.Total = result.Count()
	(&page).InitIndexes()

	var data any
	switch table {
	case "image":
		var images []model.Image
		err = result.Scan(&images)
		data = images

	case "note":
		var notes []model.Note
		err = result.Scan(&notes)
		data = notes
	}

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
