// img recycle
// @author xiangqian
// @date 13:02 2023/04/05
package recycle

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_img "note/src/api/img"
	typ_api "note/src/typ/api"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
)

// ImgList img列表
func ImgList(context *gin.Context) {
	page, err := api_img.DbPage(context, typ_api.Img{Abs: typ_api.Abs{Id: 1}})
	resp := typ_resp.Resp[typ_page.Page[typ_api.Img]]{
		Msg:  util_str.TypeToStr(err),
		Data: page,
	}
	common.HtmlOk(context, "recycle/list.html", resp)
}
