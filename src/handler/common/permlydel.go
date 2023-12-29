// @author xiangqian
// @date 22:02 2023/12/06
package common

import (
	"net/http"
	"note/src/model"
	"note/src/session"
)

func PermlyDel(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	return DelOrRestoreOrPermlyDel(request, writer, session, table, "permlydel")
}
