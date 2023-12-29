// @author xiangqian
// @date 14:13 2023/12/29
package common

import (
	"net/http"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
)

func NotFound(request *http.Request, writer http.ResponseWriter, session *session.Session, err any) (string, model.Response) {
	// 设置响应状态为 404 Not Found
	writer.WriteHeader(http.StatusNotFound)

	return "404", model.Response{Msg: util_string.String(err)}
}
