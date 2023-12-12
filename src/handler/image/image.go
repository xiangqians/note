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
	return common.List(request, writer, session, common.TableImage)
}

func Upload(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Upload(request, writer, session, common.TableImage)
}

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Rename(request, writer, session, common.TableImage)
}

func Get(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	common.Get(request, writer, session, common.TableImage)
	return "", model.Response{}
}

func View(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.View(request, writer, session, common.TableImage)
}

func Del(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Del(request, writer, session, common.TableImage)
}

func Restore(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Restore(request, writer, session, common.TableImage)
}

func PermlyDel(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.PermlyDel(request, writer, session, common.TableImage)
}
