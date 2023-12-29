// @author xiangqian
// @date 21:44 2023/12/05
package common

import (
	"net/http"
	"note/src/model"
	"note/src/session"
)

func Del(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	return DelOrRestoreOrPermlyDel(request, writer, session, table, "del")
}
