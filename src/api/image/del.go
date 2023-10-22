// @author xiangqian
// @date 23:23 2023/10/22
package image

import (
	"github.com/gin-gonic/gin"
)

// Del 图片删除
func Del(ctx *gin.Context) {
	// redirect
	//redirect := func(err any) {
	//	RedirectToList(context, err)
	//}
	//
	//// id
	//id, err := api_common_context.Param[int64](context, "id")
	//if err != nil {
	//	redirect(err)
	//	return
	//}
	//
	//// delete
	//_, err = db.Del(context, "UPDATE `lib` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", time.NowUnix(), id)
	//
	//// redirect
	//redirect(err)
}
