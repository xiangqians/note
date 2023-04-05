// recycle 回收站
// @author xiangqian
// @date 13:02 2023/04/05
package recycle

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
)

func ImgList(context *gin.Context) {
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.HtmlOk(context, "recycle/list.html", resp)
	}

	redirect(nil)
}
