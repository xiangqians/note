// @author xiangqian
// @date 21:34 2023/12/10
package note

import (
	"net/http"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
)

func List(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.List(request, writer, session, "note")
}

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Rename(request, writer, session, "note")
}
