// @author xiangqian
// @date 20:29 2023/04/27
package image

import (
	"net/http"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
)

func List(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.List(request, writer, session, "image")
}

func Upload(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Upload(request, writer, session, "image")
}

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Rename(request, writer, session, "image")
}

func Get(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	common.Get(request, writer, session, "image")
	return "", model.Response{}
}

func View(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.View(request, writer, session, "image")
}

func Del(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Del(request, writer, session, "image")
}

func Restore(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Restore(request, writer, session, "image")
}

func PermlyDel(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.PermlyDel(request, writer, session, "image")
}
