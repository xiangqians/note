// @author xiangqian
// @date 13:39 2023/12/03
package common

import (
	"fmt"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_time "note/src/util/time"
	"strconv"
	"strings"
)

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// id
	idStr := strings.TrimSpace(request.URL.Query().Get("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return redirect(table, 0, err)
	}

	// 名称
	name := strings.TrimSpace(request.URL.Query().Get("name"))
	if name == "" {
		return redirect(table, 0, err)
	}

	db := db.Get()
	sql := fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ?", table)
	result, err := db.Get(sql, id)
	if err != nil {
		return redirect(table, 0, err)
	}

	switch table {
	case "image":
		var image model.Image
		err = result.Scan(&image)
		id = image.Id
	}

	if err == nil && id != 0 {
		db.Upd(fmt.Sprintf("UPDATE `%s` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table), name, util_time.NowUnix(), id)
	}

	return redirect(table, 0, err)
}
