// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	"strings"
)

// PdfView 查看pdf文件
func PdfView(context *gin.Context, note typ_api.Note) {
	v, _ := common.Query[string](context, "v")
	v = strings.TrimSpace(v)
	switch v {
	case "1.0":
		// v1.0
	case "2.0":
		// v2.0
	default:
		v = "2.0"
	}

	note.Url = fmt.Sprintf("/note/%v", note.Id)

	resp := typ_resp.Resp[typ_api.Note]{
		Data: note,
	}
	common.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}

// DefaultView 默认查看文件
func DefaultView(context *gin.Context, note typ_api.Note, err error) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/view.html", resp)
}
