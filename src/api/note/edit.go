// note edit
// @author xiangqian
// @date 15:57 2023/04/30
package note

import (
	"github.com/gin-gonic/gin"
	typ2 "note/app/typ"
)

// Edit 文件修改页
func Edit(context *gin.Context) {
	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		FileDefaultEditPage(context, typ2.Note{}, err)
		return
	}

	// query
	f, count, err := DbQry(context, typ2.Note{Abs: typ2.Abs{Id: id}, Pid: -1})
	if err != nil || count == 0 {
		FileDefaultEditPage(context, f, err)
		return
	}

	// type
	switch typ2.ExtNameOf(f.Type) {
	// markdown
	case typ2.FtMd:
		FileMdEditPage(context, f)

	// default
	default:
		FileDefaultEditPage(context, f, err)
	}
}

// FileDefaultEditPage 默认文件修改页
func FileDefaultEditPage(context *gin.Context, note typ2.Note, err error) {
	resp := typ2.Resp[typ2.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}
	context.HtmlOk(context, "note/default/edit.html", resp)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(context *gin.Context, note typ2.Note) {
	html := func(content string, err any) {
		resp := typ2.Resp[map[string]any]{
			Msg: str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":    note,
				"content": content,
			},
		}

		context.HtmlOk(context, "note/md/edit.html", resp)
	}

	// read
	buf, err := Read(context, note)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}
