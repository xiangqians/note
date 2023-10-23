// @author xiangqian
// @date 23:21 2023/10/23
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/context"
	"note/src/dbctx"
	"note/src/model"
	util_string "note/src/util/string"
	"sort"
	"strconv"
	"strings"
)

func List[T any](ctx *gin.Context, table string) {
	html := func(page model.Page[T], err any) {
		context.HtmlOk(ctx, fmt.Sprintf("%s/list", table), model.Resp[model.Page[T]]{
			Msg:  util_string.String(err),
			Data: page,
		})
	}

	// 当前页
	current, _ := context.Query[int64](ctx, "current")
	if current <= 0 {
		current = 1
	}

	// 页数量
	size, _ := context.Query[uint8](ctx, "size")
	if size <= 0 || size > 100 {
		size = 10
	}

	var columns []string
	var values []any

	var err error

	// 要检索的字段映射，<字段名，[字段所在检索文本search中的位置，字段标识长度]>
	columnMap := map[string][]int{
		"id":   {-1, 0},
		"name": {-1, 0},
		"type": {-1, 0},
		"del":  {-1, 0},
	}
	var del byte = 0
	search, _ := context.Query[string](ctx, "search")
	if search != "" {
		cap := len(columnMap)

		// len 0, cap ?
		indexs := make([]int, 0, cap)

		// 获取字段所在检索文本search中的位置
		for name, arr := range columnMap {
			substr := fmt.Sprintf(" %s:", name)
			index := strings.Index(search, substr)
			if index == -1 {
				substr = fmt.Sprintf("%s:", name)
				index = strings.Index(search, substr)
			}

			arr[0] = index
			arr[1] = len(substr)
			indexs = append(indexs, index)
		}

		// 切片升序排序
		sort.Ints(indexs)

		// len 0, cap ?
		columns = make([]string, 0, cap)
		// len 0, cap ?
		values = make([]any, 0, cap)

		for name, arr := range columnMap {
			index := arr[0]
			if index == -1 {
				continue
			}

			var nextIndex = len(search)
			for i, index0 := range indexs {
				if index0 == index {
					i++
					if i < len(indexs) {
						nextIndex = indexs[i]
					}
					break
				}
			}

			value := search[index+arr[1] : nextIndex]
			if name == "name" {
				columns = append(columns, "`name` LIKE ?")
				value = fmt.Sprintf("%%%s%%", value)

			} else if name == "del" {
				if value == "1" {
					del = 1
				}
				continue

			} else {
				columns = append(columns, fmt.Sprintf("`%s` = ?", name))
				if name == "id" {
					var id int64
					id, err = strconv.ParseInt(value, 10, 64)
					if err != nil {
						break
					}
					value = fmt.Sprintf("%v", id)

				} else if name == "type" {
					value = strings.ToLower(value)
				}
			}
			values = append(values, value)
		}
	}

	if err != nil {
		html(model.Page[T]{
			Current: current,
			Size:    size,
			Total:   0,
			Search:  search,
		}, err)
		return
	}

	// del字段
	if columns == nil {
		// len 0, cap ?
		columns = make([]string, 0, 1)
	}
	columns = append(columns, fmt.Sprintf("`del` = %d", del))

	sql := fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `history`, `history_size`, `del`, `add_time`, `upd_time` FROM `%s`", table)
	sql += fmt.Sprintf(" WHERE %s", strings.Join(columns, " AND "))
	var page model.Page[T]
	page, err = dbctx.Page[T](ctx, current, size, sql, values...)
	page.Search = search
	html(page, err)
}

// redirectToList 重定向到图片列表
func RedirectToList(ctx *gin.Context, msg any) {
	current, _ := context.Query[int64](ctx, "current")
	size, _ := context.Query[uint8](ctx, "size")
	search, _ := context.Query[string](ctx, "search")

	// 重定向到图片首页
	context.Redirect(ctx, "/image", map[string]any{
		"current": current,
		"size":    size,
		"search":  search,
	}, msg)
}
