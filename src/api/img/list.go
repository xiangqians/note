// img list
// @author xiangqian
// @date 20:29 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	"strings"
)

// List 图片列表页面
func List(context *gin.Context) {
	// img
	img := typ_api.Img{}
	err := common.ShouldBindQuery(context, &img)

	// name
	img.Name = strings.TrimSpace(img.Name)

	// type
	ft := typ_ft.ExtNameOf(strings.TrimSpace(img.Type))
	if typ_ft.IsImg(ft) {
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
	resp := typ_resp.Resp[map[string]any]{
		Msg: util_str.TypeToStr(err),
		Data: map[string]any{
			"img":   img,   // img query
			"types": types, // types
			"page":  page,  // page
		},
	}

	// 记录查询参数
	common.SetSessionKv(context, "img", img)

	// html
	common.HtmlOk(context, "img/list.html", resp)
}
