// recycle 回收站
// @author xiangqian
// @date 13:02 2023/04/05
package recycle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// ImgRestore img恢复（还原）
func ImgRestore(context *gin.Context) {
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/recycle/img/list"), resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `del` = 0, `upd_time` = ? WHERE `id` = ?", util_time.NowUnix(), id)
	redirect(err)
	return
}

// ImgList img列表
func ImgList(context *gin.Context) {
	req, _ := common.PageReq(context)
	page, err := common.DbPage[typ_api.Img](context, req, "SELECT `id`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `img` WHERE `del` = 1 ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC")
	resp := typ_resp.Resp[typ_page.Page[typ_api.Img]]{
		Msg:  util_str.TypeToStr(err),
		Data: page,
	}
	common.HtmlOk(context, "recycle/list.html", resp)
}
