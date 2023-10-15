// 图片列表
// @author xiangqian
// @date 20:29 2023/04/27
package image

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/context"
	"note/src/typ"
	util_string "note/src/util/string"
	"strings"
)

func List(ctx *gin.Context) {
	html := func(page typ.Page[typ.Image], err any) {
		context.HtmlOk(ctx, "image/list", typ.Resp[typ.Page[typ.Image]]{
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

	search, _ := context.Query[string](ctx, "search")
	sql := "SELECT `id`, `name`, `type`, `size`, `history`, `history_size`, `del`, `add_time`, `upd_time` FROM `image` WHERE `del` = 0"
	var values []any
	if search != "" {
		arr := strings.Split(search, "&")
		for _, str := range arr {
			str = strings.TrimSpace(str)
			index := strings.Index(str, ":")
			if index > 0 && index != len(str)-1 {
				k := strings.TrimSpace(str[0:index])
				v := strings.TrimSpace(str[index+1:])
				if k == "id" {
					sql += " AND `id` = ?"
					values = append(values, v)

				} else if k == "name" {
					sql += " AND `name` LIKE ?"
					values = append(values, fmt.Sprintf("%%%s%%", v))

				} else if k == "type" {
					sql += " AND `type` = ?"
					values = append(values, v)
				}
			}
		}
	}

	// 查询
	page, err := api.DbPage[typ.Image](ctx, current, size, sql, values...)
	html(page, err)
}

//// List 库列表页面
//func List(context *gin.Context) {
//	// lib
//	lib := typ.Lib{}
//	err := api_common_context.ShouldBindQuery(context, &lib)
//
//	var History []History
//
//	// id
//	if lib.Id < 0 {
//		lib.Id = 0
//	}
//
//	// name
//	lib.Name = strings.TrimSpace(lib.Name)
//
//	// type
//	lib.Type = string(typ.ExtNameOf(strings.TrimSpace(lib.Type)))
//
//	// del
//	if lib.Del != 0 {
//		lib.Del = 1
//	}
//
//	// types
//	types := DbTypes(context)
//
//	// page
//	page, err := DbPage(context, lib)
//
//	// resp
//	resp := typ.Resp[map[string]any]{
//		Msg: str.ConvTypeToStr(err),
//		Data: map[string]any{
//			"lib":   lib,   // lib query
//			"types": types, // types
//			"page":  page,  // page
//		},
//	}
//
//	// 记录查询参数
//	session.Set(context, LibSessionKey, lib)
//
//	// html
//	api_common_context.HtmlOk(context, "lib/list.html ", resp)
//}

//// DbPage 分页查询图片
//func DbPage(context *gin.Context, img typ.Lib) (typ.Page[typ.Lib], error) {
//	// page request
//	current, size := common.PageReq(context)
//
//	// sql & args
//	args := make([]any, 0, 1)
//	sql := "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`del`, i.`add_time`, i.`upd_time` FROM `lib` i WHERE i.`del` = ? "
//	args = append(args, img.Del)
//
//	// id
//	if img.Id > 0 {
//		sql += "AND i.`id` = ? "
//		args = append(args, img.Id)
//	}
//
//	// name
//	if img.Name != "" {
//		sql += "AND i.`name` LIKE '%' || ? || '%' "
//		args = append(args, img.Name)
//	}
//
//	// type
//	if img.Type != "" {
//		sql += "AND i.`type` = ? "
//		args = append(args, img.Type)
//	}
//
//	sql += "ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC"
//
//	return db.Page[typ.Lib](context, current, size, sql, args...)
//}
//
//// DbTypes 获取库类型集合
//func DbTypes(context *gin.Context) []string {
//	types, count, err := db.Qry[[]string](context, "SELECT DISTINCT(`type`) FROM `lib` WHERE `del` = 0")
//	if err != nil || count == 0 {
//		types = nil
//	}
//
//	return types
//}