// html view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	typ2 "note/app/typ"
)

// ViewHtml 查看html文件
func ViewHtml(context *gin.Context, note typ2.Note) {
	// read
	buf, err := Read(context, note)
	if err == nil && len(buf) > 0 {
		note.Content = string(buf)
	}

	// resp
	resp := typ2.Resp[typ2.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}

	// html
	context.HtmlOk(context, "note/html/view.html", resp)
}
