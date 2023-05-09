// img restore
// @author xiangqian
// @date 21:24 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common/db"
	util_time "note/src/util/time"
)

// Restore img恢复（还原）
func Restore(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = db.DbUpd(context, "UPDATE `img` SET `del` = 0, `upd_time` = ? WHERE `del` = 1 AND `id` = ?", util_time.NowUnix(), id)
	redirect(err)
	return
}
