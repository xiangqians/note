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
)

func View(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// html模板
	html := func(data any, err any) (string, model.Response) {
		return fmt.Sprintf("%s/view", table),
			model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
				"table": table,
				"data":  data,
			}}
	}

	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return html404(err)
	}

	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return html404(err)
	}

	var data any
	switch table {
	case TableImage:
		var image model.Image
		err = result.Scan(&image)
		id = image.Id
		data = image
	}

	if err != nil || id == 0 {
		return html404(err)
	}

	return html(data, err)
}
