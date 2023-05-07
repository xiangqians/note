// html view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	util_str "note/src/util/str"
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
		Msg:  util_str.ConvTypeToStr(err),
		Data: note,
	}

	// html
	common.HtmlOk(context, "note/html/view.html", resp)
}
