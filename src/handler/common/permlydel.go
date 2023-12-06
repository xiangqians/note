// @author xiangqian
// @date 22:02 2023/12/06
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_os "note/src/util/os"
	"os"
	"strconv"
)

func PermlyDel(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return redirectList(table, err)
	}

	// 永久删除数据表记录（逻辑）
	db := db.Get()
	_, err = db.Del(fmt.Sprintf("UPDATE `%s` SET `name` = '', `type` = '', `size` = 0, `del` = 2, `add_time` = 0, `upd_time` = 0 WHERE `del` = 1 AND `id` = ?", table), id)
	if err != nil {
		return redirectList(table, err)
	}

	// 删除物理文件
	err = os.Remove(util_os.Path(dataDir, table, fmt.Sprintf("%d", id)))
	return redirectList(table, err)
}
