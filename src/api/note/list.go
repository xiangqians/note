// note list
// @author xiangqian
// @date 20:47 2023/04/09
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	"strings"
)

// List 文件列表页面
func List(context *gin.Context) {
	html := func(pnote typ_api.Note, notes []typ_api.Note, types []string, err error) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"pnote": pnote,
				"notes": notes,
				"types": types,
			},
		}
		common.HtmlOk(context, "note/list.html", resp)
	}

	// id
	pid, err := common.Query[int64](context, "pid")
	//log.Printf("id = %d\n", id)

	// name
	name, err := common.Query[string](context, "name")
	name = strings.TrimSpace(name)
	//log.Printf("name = %s\n", name)

	// type
	t, err := common.Query[string](context, "type")
	t = strings.TrimSpace(t)
	//log.Printf("t = %s\n", t)
	ft := typ_ft.ExtNameOf(t)
	if ft == typ_ft.FtUnk {
		t = ""
	}

	// pnote
	var pnote typ_api.Note
	if pid < 0 {
		pnote.Path = ""
		pnote.PathLink = ""

	} else if pid == 0 {
		pnote.Path = "/"
		pnote.PathLink = "/"

	} else {
		sql, args := DbQrySql(typ_api.Note{Abs: typ_api.Abs{Id: pid}, Pid: -1}, true)
		sql += "LIMIT 1"
		var count int64
		pnote, count, err = common.DbQry[typ_api.Note](context, sql, args...)
		if err != nil || count == 0 {
			html(pnote, nil, nil, err)
			return
		}

		if pnote.Path != "/" {
			pnote.Path += "/"
		}
		pnote.Path += fmt.Sprintf("%d:%s", pnote.Id, pnote.Name)
		InitPath(&pnote)
	}

	// list
	notes, err := DbList(context, typ_api.Note{
		Pid:  pid,
		Name: name,
		Type: t,
	})

	// types
	types, count, err := common.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	// html
	html(pnote, notes, types, err)
	return
}

func DbList(context *gin.Context, note typ_api.Note) ([]typ_api.Note, error) {
	// 查询
	path := false
	if note.Pid == -1 {
		path = true
	}
	sql, args := DbQrySql(note, path)
	sql += "LIMIT 10000"
	notes, count, err := common.DbQry[[]typ_api.Note](context, sql, args...)
	if err != nil || count == 0 {
		notes = nil
	}

	if path && err == nil && count > 0 {
		for i, l := 0, len(notes); i < l; i++ {
			InitPath(&notes[i])
		}
	}

	return notes, err
}
