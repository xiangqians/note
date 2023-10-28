// @author xiangqian
// @date 16:16 2023/10/22
package image

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/context"
	"note/src/dbctx"
	"note/src/model"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
)

// Rename 图片重命名
func Rename(ctx *gin.Context) {
	// 请求方法
	method, _ := context.Query[string](ctx, "_method")
	if method != "PUT" {
		// 重定向到图片列表
		common.RedirectToList[model.Image](ctx, "Only support request method 'PUT'")
		return
	}

	// 获取请求参数
	id, _ := context.Query[int64](ctx, "id")
	name, _ := context.Query[string](ctx, "name")

	// 校验文件名
	err := util_validate.FileName(name)

	// 文件重命名
	if err == nil && id > 0 {
		_, err = dbctx.Upd(ctx, "UPDATE `image` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", name, util_time.NowUnix(), id)
	}

	// 重定向到图片列表
	common.RedirectToList[model.Image](ctx, err)
}
