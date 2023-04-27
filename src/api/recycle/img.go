// img recycle
// @author xiangqian
// @date 13:02 2023/04/05
package recycle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_img "note/src/api/img"
	typ_api "note/src/typ/api"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// ImgPermlyDel img永久删除
func ImgPermlyDel(context *gin.Context) {
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

	// img
	img, count, err := api_img.DbQry(context, id, 1)
	if err != nil || count == 0 {
		redirect(err)
		return
	}

	// 删除图片历史记录
	hist := img.Hist
	if hist != "" {
		hists := make([]typ_api.Img, 0, 1) // len 0, cap ?
		err = util_json.Deserialize(hist, &hists)
		if err != nil {
			redirect(err)
			return
		}

		for _, histImg := range hists {
			path, err := api_img.HistPath(context, histImg)
			if err == nil {
				util_os.DelFile(path)
			}
		}
	}

	// 删除图片
	path, err := api_img.Path(context, img)
	if err == nil {
		util_os.DelFile(path)
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `img` SET `name` = '', `type` = '', `size` = 0, `hist` = '', `hist_size` = 0, `del` = 2, `add_time` = 0, `upd_time` = 0 WHERE `id` = ?", id)
	redirect(err)
	return
}

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
	page, err := api_img.DbPage(context, typ_api.Img{Abs: typ_api.Abs{Id: 1}})
	resp := typ_resp.Resp[typ_page.Page[typ_api.Img]]{
		Msg:  util_str.TypeToStr(err),
		Data: page,
	}
	common.HtmlOk(context, "recycle/list.html", resp)
}
