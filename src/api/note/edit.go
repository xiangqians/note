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
		FileDefaultEditPage(context, typ.Note{}, err)
		return
	}

	// query
	f, count, err := DbQry(context, typ.Note{Abs: typ.Abs{Id: id}, Pid: -1})
	if err != nil || count == 0 {
		FileDefaultEditPage(context, f, err)
		return
	}

	// type
	switch typ.ExtNameOf(f.Type) {
	// markdown
	case typ.FtMd:
		FileMdEditPage(context, f)

	// default
	default:
		FileDefaultEditPage(context, f, err)
	}
}

// FileDefaultEditPage 默认文件修改页
func FileDefaultEditPage(context *gin.Context, note typ.Note, err error) {
	resp := typ.Resp[typ.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}
	api_common_context.HtmlOk(context, "note/default/edit.html", resp)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(context *gin.Context, note typ.Note) {
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

	html(content, err)
}
