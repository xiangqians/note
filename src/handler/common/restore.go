// @author xiangqian
// @date 22:09 2023/12/06
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

func Restore(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return RedirectList(table, nil, err)
	}

	db := db.Get()
	_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `del` = 0, `upd_time` = ? WHERE `del` = 1 AND `id` = ?", table), util_time.NowUnix(), id)
	return RedirectList(table, nil, err)
}
