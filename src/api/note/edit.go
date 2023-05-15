// note edit
// @author xiangqian
// @date 15:57 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
)

// Edit 文件修改页
func Edit(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		EditUnsupported(context, typ.Note{}, err)
		return
	}

	// query
	f, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 {
		EditUnsupported(context, f, err)
		return
	}

	// type
	switch typ.ExtNameOf(f.Type) {

	// markdown
	case typ.FtMd:
		EditMd(context, f)

	// unsupported
	default:
		EditUnsupported(context, f, err)
	}
}

// EditUnsupported 不支持编辑
func EditUnsupported(context *gin.Context, note typ.Note, err any) {
	resp := typ.Resp[typ.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}
	api_common_context.HtmlOk(context, "note/unsupported/edit.html", resp)
}
