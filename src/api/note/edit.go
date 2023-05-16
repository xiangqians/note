// note edit
// @author xiangqian
// @date 15:57 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
)

// Edit 文件修改页
func Edit(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		common.DataNotExist(context, err)
		return
	}

	// query
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 {
		common.DataNotExist(context, err)
		return
	}

	// type
	switch typ.ExtNameOf(note.Type) {
	// markdown
	case typ.FtMd:
		EditMd(context, note)

	// unsupported，不支持编辑
	default:
		resp := typ.Resp[typ.Note]{
			Msg:  str.ConvTypeToStr(err),
			Data: note,
		}
		api_common_context.HtmlOk(context, "note/unsupported/edit.html", resp)
	}
}
