// html view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
)

// HtmlView 查看html文件
func HtmlView(context *gin.Context, note typ_api.Note) {
	html := func(html string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note": note,
				"html": html,
			},
		}
		common.HtmlOk(context, "note/html/view.html", resp)
	}

	// read
	buf, err := Read(context, note)
	if err != nil {
		html("", err)
		return
	}

	html(string(buf), nil)
}
