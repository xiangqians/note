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
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return redirectList(table, nil, err)
	}

	db := db.Get()
	_, err = db.Del(fmt.Sprintf("UPDATE `%s` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table), util_time.NowUnix(), id)
	return redirectList(table, nil, err)
}
