// note list
// @author xiangqian
// @date 20:47 2023/04/09
package note

import (
	"github.com/gin-gonic/gin"
	"note/app/api/common/db"
	"note/app/api/common/session"
	typ2 "note/app/typ"
	"strings"
)

// List 文件列表页面
func List(context *gin.Context) {
	html := func(note typ2.Note, types []string, err error) {
		resp := typ2.Resp[map[string]any]{
			Msg: str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":  note,
				"types": types,
			},
		}
		context.HtmlOk(context, "note/list.html", resp)
	}

	// note
	note := typ2.Note{}
	err := context.ShouldBindQuery(context, &note)
	note.Del = 0
	note.QryPath = 0
	note.Children = nil
	//if err != nil {
	//	html(note, nil, err)
	//	return
	//}

	// name
	note.Name = strings.TrimSpace(note.Name)
	//log.Printf("name = %s\n", name)

	// type
	note.Type = strings.TrimSpace(note.Type)

	// p
	pid := note.Pid
	if pid < 0 {
		html(note, nil, err)
		return
	}

	var pNote typ2.Note
	if pid != 0 {
		var count int64
		pNote, count, err = DbQry(context, typ2.Note{Abs: typ2.Abs{Id: pid}, Pid: -1, QryPath: 2})
		if err != nil || count == 0 {
			html(note, nil, err)
			return
		}
	}

	// types
	types, count, err := db.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	// list
	if note.Sub != 0 {
		note.QryPath = 1
	}
	children, count, err := DbList(context, note)
	if err == nil && count > 0 {
		note.Children = children
	}

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
	session.SetSessionKv(context, "note", typ2.Note{
		Pid:     note.Pid,
		Deleted: note.Deleted,
	})

	// html
	html(note, types, err)
	return
}
