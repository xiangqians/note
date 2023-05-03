// note pdf view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	"strings"
)

// ViewPdf 查看pdf文件
func ViewPdf(context *gin.Context, note typ_api.Note) {
	v, _ := common.Query[string](context, "v")
	v = strings.TrimSpace(v)
	switch v {
	// v1.0
	case "1.0":

	// v2.0
	case "2.0":

	// default v2.0
	default:
		v = "2.0"
	}

	note.Url = fmt.Sprintf("/note/%v", note.Id)

	resp := typ_resp.Resp[typ_api.Note]{
		Data: note,
	}
	common.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}
