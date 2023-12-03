// @author xiangqian
// @date 20:29 2023/04/27
package image

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
)

func List(request *http.Request, session *session.Session) (string, model.Response) {
	return common.List(request, session, "image")
}

func Rename(request *http.Request, session *session.Session) (string, model.Response) {
	return common.Rename(request, session, "image")
}

func Rename1(request *http.Request, session *session.Session) (string, model.Response) {
	vars := mux.Vars(request)
	id := vars["id"]
	log.Println("------------2: ", id, request.URL.Path)
	return common.List(request, session, "image")
}
