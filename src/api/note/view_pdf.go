// note pdf view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	typ2 "note/app/typ"
	"strings"
)

// ViewPdf 查看pdf文件
func ViewPdf(context *gin.Context, note typ2.Note) {
	// version
	v, _ := context.Query[string](context, "v")
	v = strings.TrimSpace(v)
	switch v {
	// v1.0
	case "1.0":

	// v2.0
	case "2.0":

	// default v2.0
	default:
		v = "2.0"
	}

	// resp
	resp := typ2.Resp[typ2.Note]{
		Data: note,
	}

	// html
	context.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}
