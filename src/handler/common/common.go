// @author xiangqian
// @date 19:45 2023/12/04
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_i18n "note/src/util/i18n"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	"os"
	"strconv"
	"strings"
)

// DataDir 数据目录
var DataDir = model.Ini.Data.Dir

// RedirectList 重定向到列表
func RedirectList(table string, paramMap map[string]any, err any) (string, model.Response) {
	name := fmt.Sprintf("redirect:/%s", table)
	if paramMap != nil {
		// len 0, cap ?
		params := make([]string, 0, len(paramMap))
		for key, value := range paramMap {
			params = append(params, fmt.Sprintf("%s=%s", key, url.QueryEscape(util_string.String(value))))
		}
		name += "?" + strings.Join(params, "&")
	}
	return name, model.Response{Msg: util_string.String(err)}
}

// RedirectView 重定向到详情页
func RedirectView(table string, id int64, err any) (string, model.Response) {
	return fmt.Sprintf("redirect:/%s/%d/view", table, id), model.Response{Msg: util_string.String(err)}
}

func DelOrRestoreOrPermlyDel(request *http.Request, writer http.ResponseWriter, session *session.Session, table string, Type string) (string, model.Response) {
	var pid int64 = 0
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

	db := db.Get()

	// note.pid
	if table == TableNote {
		result, err := db.Get(fmt.Sprintf("SELECT `pid` FROM `%s` WHERE `del` IN (0, 1) AND `id` = ?", table), id)
		if err != nil {
			return redirect(err)
		}

		err = result.Scan(&pid)
		if err != nil {
			return redirect(err)
		}
	}

	switch Type {
	// 删除（逻辑）
	case "del":
		if table == TableNote {
			result, err := db.Get("SELECT COUNT(1) FROM `note` WHERE `del` IN (0, 1) AND `pid` = ?", id)
			if err != nil {
				return redirect(err)
			}

			var count int64
			err = result.Scan(&count)
			if err != nil {
				return redirect(err)
			}

			if count != 0 {
				return redirect(util_i18n.GetMessage("i18n.cannotDelNonEmptyFolder", session.GetLanguage()))
			}
		}
		_, err = db.Del(fmt.Sprintf("UPDATE `%s` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table), util_time.NowUnix(), id)

	// 恢复
	case "restore":
		_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `del` = 0, `upd_time` = ? WHERE `del` = 1 AND `id` = ?", table), util_time.NowUnix(), id)

	// 永久删除数据表记录（逻辑）
	case "permlydel":
		_, err = db.Del(fmt.Sprintf("UPDATE `%s` SET `name` = '', `type` = '', `size` = 0, `del` = 2, `add_time` = 0, `upd_time` = 0 WHERE `del` = 1 AND `id` = ?", table), id)
		if err == nil {
			// 删除物理文件
			err = os.Remove(util_os.Path(DataDir, table, fmt.Sprintf("%d", id)))
			log.Println(err)
			err = nil
		}
	}
	return redirect(err)
}

func getPNoteSql() string {
	return "SELECT" +
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
		" FROM `note` p1" +
		" LEFT JOIN `note` p2 ON p2.`id` = p1.`pid`" +
		" LEFT JOIN `note` p3 ON p3.`id` = p2.`pid`" +
		" LEFT JOIN `note` p4 ON p4.`id` = p3.`pid`" +
		" LEFT JOIN `note` p5 ON p5.`id` = p4.`pid`" +
		" LEFT JOIN `note` p6 ON p6.`id` = p5.`pid`" +
		" LEFT JOIN `note` p7 ON p7.`id` = p6.`pid`" +
		" LEFT JOIN `note` p8 ON p8.`id` = p7.`pid`" +
		" LEFT JOIN `note` p9 ON p9.`id` = p8.`pid`" +
		" LEFT JOIN `note` p10 ON p10.`id` = p9.`pid`" +
		" WHERE p1.`id` = ?"
}
