// @author xiangqian
// @date 13:41 2023/11/19
package index

import (
	"net/http"
	"note/src/model"
	"note/src/session"
)

func Index(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return "index", model.Response{}
}
