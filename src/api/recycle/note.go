// note recycle
// @author xiangqian
// @date 14:16 2023/04/08
package recycle

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_note "note/src/api/note"
	typ_api "note/src/typ/api"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
)

// NoteList note列表
func NoteList(context *gin.Context) {
	page, err := api_note.DbPage(context, 1)
	resp := typ_resp.Resp[typ_page.Page[typ_api.Note]]{
		Msg:  util_str.TypeToStr(err),
		Data: page,
	}
	common.HtmlOk(context, "recycle/list.html", resp)
}
