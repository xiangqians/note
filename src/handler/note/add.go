// @author xiangqian
// @date 22:44 2023/12/11
package note

import (
	"fmt"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	"note/src/util/time"
	util_validate "note/src/util/validate"
	"strconv"
	"strings"
)

func add(request *http.Request, writer http.ResponseWriter, session *session.Session, Type string) (string, model.Response) {
	// 父id
	pid, err := strconv.ParseInt(strings.TrimSpace(request.PostFormValue("pid")), 10, 64)

	redirect := func(err any) (string, model.Response) {
		return "redirect:/note?search=pid%3A%20" + fmt.Sprintf("%d", pid), model.Response{}
	}

	if pid < 0 {
		return redirect(nil)
	}

	// 名称
	name := strings.TrimSpace(request.PostFormValue("name"))
	if name == "" {
		return redirect(nil)
	}
	err = util_validate.FileName(name, session.GetLanguage())
	if err != nil {
		return redirect(nil)
	}

	db := db.Get()

	// 校验父id是否存在
	result, err := db.Get("SELECT `id` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", pid)
	if err != nil {
		return redirect(err)
	}
	pid = 0
	err = result.Scan(&pid)
	if err != nil {
		return redirect(err)
	}
	if pid == 0 {
		return redirect(nil)
	}

	// 获取永久删除id，以复用
	result, err = db.Get("SELECT `id` FROM `note` WHERE `del` = 2 LIMIT 1")
	if err != nil {
		return redirect(err)
	}
	var id int64
	err = result.Scan(&id)
	if err != nil {
		return redirect(err)
	}

	// 新id
	if id == 0 {
		_, id, err = db.Add("INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", pid, name, Type, time.NowUnix())
	} else
	// 复用id
	{
		_, err = db.Upd("UPDATE `note` SET `pid` = ?, `name` = ?, `type` = ?, `size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?",
			pid,
			name,
			Type,
			time.NowUnix(),
			id)
	}

	return redirect(err)
}
