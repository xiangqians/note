// edit md
// @author xiangqian
// @date 15:38 2023/05/13
package note

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
)

// EditMd md文件修改页
func EditMd(context *gin.Context, note typ.Note) {
	html := func(content string, err any) {
		resp := typ.Resp[map[string]any]{
			Msg: str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":    note,
				"content": content,
			},
		}

		api_common_context.HtmlOk(context, "note/md/edit.html", resp)
	}

	// read
	buf, err := Read(context, note)
	content := ""
	if err == nil {
		content = string(buf)
	}

	// html
	html(content, err)
}
