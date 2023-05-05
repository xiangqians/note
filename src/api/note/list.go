// note list
// @author xiangqian
// @date 20:47 2023/04/09
package note

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	"strings"
)

// List 文件列表页面
func List(context *gin.Context) {
	html := func(note typ_api.Note, types []string, err error) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":  note,
				"types": types,
			},
		}
		common.HtmlOk(context, "note/list.html", resp)
	}

	// note
	note := typ_api.Note{}
	err := common.ShouldBindQuery(context, &note)
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

	var pNote typ_api.Note
	if pid != 0 {
		var count int64
		pNote, count, err = DbQry(context, typ_api.Note{Abs: typ_api.Abs{Id: pid}, Pid: -1, QryPath: 2})
		if err != nil || count == 0 {
			html(note, nil, err)
			return
		}
	}

	// types
	types, count, err := common.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
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
	common.SetSessionKv(context, "note", typ_api.Note{
		Pid:     note.Pid,
		Deleted: note.Deleted,
	})

	// html
	html(note, types, err)
	return
}
