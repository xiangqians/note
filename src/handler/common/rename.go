// @author xiangqian
// @date 13:39 2023/12/03
package common

import (
	"fmt"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"strconv"
	"strings"
)

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	redirect := func(pid int64, err any) (string, model.Response) {
		var name string
		if table == TableNote {
			name = fmt.Sprintf("redirect:/%s/%d/list", table, pid)
		} else {
			name = fmt.Sprintf("redirect:/%s", table)
		}
		return name, model.Response{Msg: util_string.String(err)}
	}

	var pid int64

	// id
	id, err := strconv.ParseInt(strings.TrimSpace(request.PostFormValue("id")), 10, 64)
	if err != nil || id <= 0 {
		return redirect(pid, err)
	}

	var result *db.Result
	db := db.Get()

	// 校验id
	if table == TableNote {
		result, err = db.Get(fmt.Sprintf("SELECT `id`, `pid` FROM `%s` WHERE `del` = 0 AND `id` = ?", table), id)
		if err != nil {
			return redirect(pid, err)
		}
		var note model.Note
		err = result.Scan(&note)
		id = note.Id
		pid = note.Pid

	} else {
		result, err = db.Get(fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 0 AND `id` = ?", table), id)
		if err != nil {
			return redirect(pid, err)
		}
		id = 0
		err = result.Scan(&id)
	}
	if err != nil || id <= 0 {
		return redirect(pid, err)
	}

	// 名称
	name := strings.TrimSpace(request.PostFormValue("name"))
	if name == "" {
		return redirect(pid, err)
	}
	err = util_validate.FileName(name, session.GetLanguage())
	if err != nil {
		return redirect(pid, err)
	}

	// 更新名称
	_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` != ?", table), name, util_time.NowUnix(), id, name)
	return redirect(pid, err)
}
