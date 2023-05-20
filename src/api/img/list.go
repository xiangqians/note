// img list
// @author xiangqian
// @date 20:29 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/str"
	"strings"
)

// List 图片列表页面
func List(context *gin.Context) {
	// img
	img := typ.Img{}
	err := api_common_context.ShouldBindQuery(context, &img)

	// name
	img.Name = strings.TrimSpace(img.Name)

	// type
	img.Type = string(typ.ExtNameOf(strings.TrimSpace(img.Type)))

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
		Msg: str.ConvTypeToStr(err),
		Data: map[string]any{
			"img":   img,   // img query
			"types": types, // types
			"page":  page,  // page
		},
	}

	// 记录查询参数
	session.Set(context, ImgSessionKey, img)

	// html
	api_common_context.HtmlOk(context, "img/list.html", resp)
}

// DbPage 分页查询图片
func DbPage(context *gin.Context, img typ.Img) (typ.Page[typ.Img], error) {
	// page request
	current, size := common.PageReq(context)

	// sql & args
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

	return db.Page[typ.Img](context, current, size, sql, args...)
}

// DbTypes 获取图片类型集合
func DbTypes(context *gin.Context) []string {
	types, count, err := db.Qry[[]string](context, "SELECT DISTINCT(`type`) FROM `img` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	return types
}
