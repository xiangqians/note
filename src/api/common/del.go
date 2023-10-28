// @author xiangqian
// @date 15:11 2023/10/28
package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/context"
	"note/src/dbctx"
	"note/src/util/time"
)

// Del 删除
func Del[T any](ctx *gin.Context) {
	// 请求方法
	requestMethod := RequestMethod(ctx)
	if requestMethod != http.MethodDelete {
		// 重定向到列表
		RedirectToList[T](ctx, "Only support request method 'DELETE'")
		return
	}

	// id
	id, err := context.Param[int64](ctx, "id")
	if err != nil || id <= 0 {
		// 重定向到列表
		RedirectToList[T](ctx, err)
		return
	}

	// 数据表名
	table := Table[T]()

	// delete
	_, err = dbctx.Del(ctx, "UPDATE `%s` SET `del` = 1, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table, time.NowUnix(), id)

	// 重定向到列表
	RedirectToList[T](ctx, err)
}

// PermlyDel 永久删除
func PermlyDel(ctx *gin.Context) {

}
