// index
// @author xiangqian
// @date 21:38 2023/07/11
package index

import (
	"github.com/gin-gonic/gin"
	"note/src/context"
	"note/src/typ"
)

// Index index页面
func Index(ctx *gin.Context) {
	context.HtmlOk(ctx, "index", typ.Resp[any]{})
}
