// img del
// @author xiangqian
// @date 21:48 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/util/os"
	"note/src/util/time"
)

// PermlyDel 永久删除图片
func PermlyDel(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := api_common_context.Param[int64](context, "id")
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
	histImgs, err := DeserializeHist(img.Hist)
	if err != nil {
		redirect(err)
		return
	}
	if histImgs != nil {
		for _, histImg := range histImgs {
			path, err := HistPath(context, histImg)
			if err == nil {
				os.DelFile(path)
			}
		}
	}

	// 删除图片
	path, err := Path(context, img)
	if err == nil {
		os.DelFile(path)
	}

	// delete
	_, err = db.Del(context, "UPDATE `img` SET `name` = '', `type` = '', `size` = 0, `hist` = '', `hist_size` = 0, `del` = 2, `add_time` = 0, `upd_time` = 0 WHERE `del` = 1 AND `id` = ?", id)

	// redirect
	redirect(err)
}

// Del 删除图片
func Del(context *gin.Context) {
	// redirect
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// delete
	_, err = db.Del(context, "UPDATE `img` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", time.NowUnix(), id)

	// redirect
	redirect(err)
}
