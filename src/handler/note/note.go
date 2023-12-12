// @author xiangqian
// @date 21:34 2023/12/10
package note

import (
	"net/http"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
	util_filetype "note/src/util/filetype"
)

func List(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.List(request, writer, session, common.TableNote)
}

func Rename(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Rename(request, writer, session, common.TableNote)
}

func AddFolder(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return add(request, writer, session, util_filetype.Folder)
}

func AddMdFile(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return add(request, writer, session, util_filetype.Md)
}

func Upload(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return common.Upload(request, writer, session, common.TableNote)
}
