// img restore
// @author xiangqian
// @date 21:24 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/util/time"
)

// Restore img恢复（还原）
func Restore(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = db.Upd(context, "UPDATE `img` SET `del` = 0, `upd_time` = ? WHERE `del` = 1 AND `id` = ?", time.NowUnix(), id)

	// redirect
	redirect(err)
}
