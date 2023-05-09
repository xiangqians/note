// img list
// @author xiangqian
// @date 20:29 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	util_str "note/src/util/str"
	"strings"
)

// List 图片列表页面
func List(context *gin.Context) {
	// img
	img := typ.Img{}
	err := context.ShouldBindQuery(context, &img)

	// name
	img.Name = strings.TrimSpace(img.Name)

	// type
	ft := typ.ExtNameOf(strings.TrimSpace(img.Type))
	if typ.IsImg(ft) {
		img.Type = string(ft)
	} else {
		img.Type = ""
	}

	// del
	if !(img.Del == 0 || img.Del == 1) {
		img.Del = 0
	}

	// types
	types := DbTypes(context)

	// page
	page, err := DbPage(context, img)

	// resp
	resp := typ.Resp[map[string]any]{
		Msg: util_str.ConvTypeToStr(err),
		Data: map[string]any{
			"img":   img,   // img query
			"types": types, // types
			"page":  page,  // page
		},
	}

	// 记录查询参数
	session.SetSessionKv(context, "img", img)

	// html
	context.HtmlOk(context, "img/list.html", resp)
}

// DbPage 分页查询图片
func DbPage(context *gin.Context, img typ.Img) (typ.Page[typ.Img], error) {
	current, _ := context.Param[int64](context, "current")
	if current <= 0 {
		current = 1
	}

	size, _ := context.Param[uint8](context, "size")
	if size <= 0 {
		size = 10
	}

	args := make([]any, 0, 1)
	sql := "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`del`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = ? "
	args = append(args, img.Del)

	// id
	if img.Id > 0 {
		sql += "AND i.`id` = ? "
		args = append(args, img.Id)
	}

	// name
	if img.Name != "" {
		sql += "AND i.`name` LIKE '%' || ? || '%' "
		args = append(args, img.Name)
	}

	// type
	if img.Type != "" {
		sql += "AND i.`type` = ? "
		args = append(args, img.Type)
	}

	sql += "ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC"

	return db.DbPage[typ.Img](context, current, size, sql, args...)
}

// DbTypes 获取图片类型集合
func DbTypes(context *gin.Context) []string {
	// types
	types, count, err := db.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `img` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	return types
}
