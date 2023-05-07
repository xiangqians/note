// note pdf view
// @author xiangqian
// @date 15:56 2023/04/30
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	"strings"
)

// ViewPdf 查看pdf文件
func ViewPdf(context *gin.Context, note typ.Note) {
	// version
	v, _ := common.Query[string](context, "v")
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
	resp := typ.Resp[typ.Note]{
		Data: note,
	}

	// html
	common.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}
