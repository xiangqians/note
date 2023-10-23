// @author xiangqian
// @date 23:23 2023/10/22
package image

import (
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/context"
	"note/src/util/time"
)

// Del 图片删除
func Del(ctx *gin.Context) {
	// 请求方法
	method, _ := context.Query[string](ctx, "_method")
	if method != "DELETE" {
		// 重定向到图片列表
		redirectToList(ctx, "Only support request method 'DELETE'")
		return
	}

	// id
	id, err := context.Param[int64](ctx, "id")
	if err != nil || id <= 0 {
		// 重定向到图片列表
		redirectToList(ctx, err)
		return
	}

	// delete
	_, err = api.DbExec(ctx, "UPDATE `image` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", time.NowUnix(), id)

	// 重定向到图片列表
	redirectToList(ctx, err)
}
