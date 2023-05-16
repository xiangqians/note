// note pdf view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"strings"
)

// ViewPdf 查看pdf文件
func ViewPdf(context *gin.Context, note typ.Note) {
	// version
	v, _ := api_common_context.Query[string](context, "v")
	v = strings.TrimSpace(v)
	if !(v == "1.0" || v == "2.0") {
		//  default v2.0
		v = "2.0"
	}

	// resp
	resp := typ.Resp[typ.Note]{
		Data: note,
	}

	// html
	api_common_context.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}
