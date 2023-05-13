// html view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
)

// ViewHtml 查看html文件
func ViewHtml(context *gin.Context, note typ.Note) {
	// read
	buf, err := Read(context, note)
	if err == nil && len(buf) > 0 {
		note.Content = string(buf)
	}

	// resp
	resp := typ.Resp[typ.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}

	// html
	api_common_context.HtmlOk(context, "note/html/view.html", resp)
}
