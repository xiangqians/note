// index
// @author xiangqian
// @date 21:38 2023/07/11
package index

import (
	"github.com/gin-gonic/gin"
	src_context "note/src/context"
	"note/src/typ"
)

// Index index页面
func Index(context *gin.Context) {
	src_context.HtmlOk(context, "index", typ.Resp[any]{})
}
