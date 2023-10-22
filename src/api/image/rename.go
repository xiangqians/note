// 图片重命名
// @author xiangqian
// @date 16:16 2023/10/22
package image

import (
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/context"
	"note/src/session"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
)

func Rename(ctx *gin.Context) {
	id, _ := context.Query[int64](ctx, "id")
	name, _ := context.Query[string](ctx, "name")
	current, _ := context.Query[int64](ctx, "current")
	size, _ := context.Query[uint8](ctx, "size")
	search, _ := context.Query[string](ctx, "search")

	// 校验文件名
	err := util_validate.FileName(name)

	// 文件重命名
	if id > 0 {
		_, err = api.DbExec(ctx, "UPDATE `image` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", name, util_time.NowUnix(), id)
	}

	// 判断是否有错误，有错误则存储到session中
	if err != nil {
		session.Set(ctx, imageErrKey, util_string.String(err))
	}

	// 重定向到图片首页
	context.Redirect(ctx, "/image", map[string]any{
		"current": current,
		"size":    size,
		"search":  search,
	})
}
