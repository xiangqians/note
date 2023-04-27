// img del
// @author xiangqian
// @date 21:48 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
)

// PermlyDel 永久删除图片
func PermlyDel(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 1)
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
			path, err := HistPath(context, histImg)
			if err == nil {
				util_os.DelFile(path)
			}
		}
	}

	// 删除图片
	path, err := Path(context, img)
	if err == nil {
		util_os.DelFile(path)
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `img` SET `name` = '', `type` = '', `size` = 0, `hist` = '', `hist_size` = 0, `del` = 2, `add_time` = 0, `upd_time` = 0 WHERE `del` = 1 AND `id` = ?", id)
	redirect(err)
	return
}

// Del 删除图片
func Del(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `img` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", util_time.NowUnix(), id)
	redirect(err)
	return
}
