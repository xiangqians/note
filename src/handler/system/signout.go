// @author xiangqian
// @date 21:41 2023/07/11
package system

import (
	"net/http"
	"note/src/model"
	"note/src/session"
)

// SignOut 注销
func SignOut(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	// 清除会话
	session.Clear()

	// 重定向到登录页
	return "redirect:/signin", model.Response{}
}
