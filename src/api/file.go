// file
// @author xiangqian
// @date 17:50 2023/02/04
package api

import "github.com/gin-gonic/gin"

func FileListPage(pContext *gin.Context) {
	Html(pContext, "file/list.html", nil, nil)
	return
}
