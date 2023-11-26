// @author xiangqian
// @date 20:29 2023/04/27
package image

import (
	"net/http"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
)

func List(request *http.Request, session *session.Session) (string, model.Response) {
	return common.List(request, session, "image")
}
