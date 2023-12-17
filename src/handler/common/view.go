// @author xiangqian
// @date 22:39 2023/12/04
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
	"strconv"
	"strings"
)

func View(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return NotFound(err)
	}

	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return NotFound(err)
	}

	var abs model.Abs
	err = result.Scan(&abs)
	if err != nil || abs.Id == 0 {
		return NotFound(err)
	}

	if table == TableNote {
		sql := "SELECT" +
			"    (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`id` END)" +
			" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`id`END)" +
			" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`id`END)" +
			" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`id`END)" +
			" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`id`END)" +
			" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`id`END)" +
			" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`id`END)" +
			" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`id`END)" +
			" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`id`END)" +
			" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`id`END) AS 'ids_str'," +
			"    (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`name` END)" +
			" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`name` END)" +
			" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`name` END)" +
			" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`name` END)" +
			" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`name` END)" +
			" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`name` END)" +
			" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`name` END)" +
			" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`name` END)" +
			" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`name` END)" +
			" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`name` END) AS 'names_str'" +
			" FROM `note` t" +
			" LEFT JOIN `note` p1 ON p1.`id` = t.`pid`" +
			" LEFT JOIN `note` p2 ON p2.`id` = p1.`pid`" +
			" LEFT JOIN `note` p3 ON p3.`id` = p2.`pid`" +
			" LEFT JOIN `note` p4 ON p4.`id` = p3.`pid`" +
			" LEFT JOIN `note` p5 ON p5.`id` = p4.`pid`" +
			" LEFT JOIN `note` p6 ON p6.`id` = p5.`pid`" +
			" LEFT JOIN `note` p7 ON p7.`id` = p6.`pid`" +
			" LEFT JOIN `note` p8 ON p8.`id` = p7.`pid`" +
			" LEFT JOIN `note` p9 ON p9.`id` = p8.`pid`" +
			" LEFT JOIN `note` p10 ON p10.`id` = p9.`pid`" +
			" WHERE t.`id` = ?"
		result, err = db.Get(sql, id)
		if err != nil {
			return NotFound(err)
		}

		var pNote model.PNote
		err = result.Scan(&pNote)
		if err != nil {
			return NotFound(err)
		}

		if pNote.IdsStr != "" {
			pNote.Ids = strings.Split(pNote.IdsStr, "/")[1:]
			pNote.Names = strings.Split(pNote.NamesStr, "/")[1:]
		}

		return fmt.Sprintf("%s/view", table),
			model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
				"table": table,
				"data":  abs,
				"pNote": pNote,
			}}
	}

	return fmt.Sprintf("%s/view", table),
		model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
			"table": table,
			"data":  abs,
		}}
}
