// note list
// @author xiangqian
// @date 20:47 2023/04/09
package note

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/str"
	"strings"
)

// List 文件列表页面
func List(context *gin.Context) {
	html := func(note typ.Note, types []string, err error) {
		resp := typ.Resp[map[string]any]{
			Msg: str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":  note,
				"types": types,
			},
		}
		api_common_context.HtmlOk(context, "note/list.html", resp)
	}

	// note
	note := typ.Note{}
	err := api_common_context.ShouldBindQuery(context, &note)

	// name
	note.Name = strings.TrimSpace(note.Name)

	// type
	note.Type = string(typ.ExtNameOf(strings.TrimSpace(note.Type)))

	// pid
	pid := note.Pid
	if pid < 0 {
		html(note, nil, err)
		return
	}

	// p note
	var pNote typ.Note
	if pid != 0 {
		var count int64
		pNote, count, err = DbQry(context, pid, 2, 0)
		if err != nil || count == 0 {
			html(note, nil, err)
			return
		}
	}

	// types
	types, count, err := db.Qry[[]string](context, "SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	// list
	if note.ContainsSub != 0 {
		note.QryPath = 1
	}
	children, count, err := DbList(context, note)
	if err != nil || count == 0 {
		children = nil
	}
	note.Children = children

	if pid == 0 {
		note.Id = 0
		note.Pid = 0
		note.Path = "/"
		note.PathLink = "/"
	} else {
		note.Id = pNote.Id
		note.Pid = pNote.Pid
		note.Path = pNote.Path
		note.PathLink = pNote.PathLink
	}

	// 记录查询参数
	session.Set(context, NoteSessionKey, typ.Note{
		Abs: typ.Abs{Del: note.Del},
		Pid: note.Pid,
	})

	// html
	html(note, types, err)
}
