// @author xiangqian
// @date 13:41 2023/11/19
package index

import (
	"net/http"
	"note/src/model"
)

func Index(r *http.Request) (string, model.Response) {
	return "index.html", model.Response{}
}
