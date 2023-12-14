// @author xiangqian
// @date 21:44 2023/12/05
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_time "note/src/util/time"
	"strconv"
)

func Del(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	var pid int64 = 0
	redirect := func(err any) (string, model.Response) {
		var paramMap map[string]any = nil
		if table == TableNote {
			paramMap = map[string]any{"search": fmt.Sprintf("pid: %d", pid)}
		}
		return RedirectList(table, paramMap, err)
	}

	// id
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return redirect(err)
	}

	db := db.Get()

	// 获取pid
	if table == TableNote {
		result, err := db.Get(fmt.Sprintf("SELECT `pid` FROM `%s` WHERE `del` = 0 AND `id` = ?", table), id)
		if err != nil {
			return redirect(err)
		}
		err = result.Scan(&pid)
		if err != nil {
			return redirect(err)
		}
	}

	_, err = db.Del(fmt.Sprintf("UPDATE `%s` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table), util_time.NowUnix(), id)
	return redirect(err)
}
