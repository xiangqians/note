// @author xiangqian
// @date 13:39 2023/12/03
package common

import (
	"fmt"
	"github.com/gorilla/mux"
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
	var pid int64
	redirect := func(err any) (string, model.Response) {
		if table == TableNote {
			return fmt.Sprintf("redirect:/%s/%d/list", table, pid), model.Response{Msg: util_string.String(err)}
		}
		return fmt.Sprintf("redirect:/%s", table), model.Response{Msg: util_string.String(err)}
	}

	// id
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return redirect(err)
	}

	var result *db.Result
	db := db.Get()

	// 校验id
	// note
	if table == TableNote {
		result, err = db.Get(fmt.Sprintf("SELECT `id`, `pid` FROM `%s` WHERE `del` = 0 AND `id` = ?", table), id)
		if err != nil {
			return redirect(err)
		}
		var note model.Note
		err = result.Scan(&note)
		id = note.Id
		pid = note.Pid
	} else
	// abs
	{
		result, err = db.Get(fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 0 AND `id` = ?", table), id)
		if err != nil {
			return redirect(err)
		}
		id = 0
		err = result.Scan(&id)
	}
	if err != nil || id == 0 {
		return redirect(err)
	}

	// 名称
	name := strings.TrimSpace(request.PostFormValue("name"))
	if name == "" {
		return redirect(err)
	}
	err = util_validate.FileName(name, session.GetLanguage())
	if err != nil {
		return redirect(err)
	}

	// 更新名称
	_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` != ?", table), name, util_time.NowUnix(), id, name)
	return redirect(err)
}
