// lib update
// @author xiangqian
// @date 21:58 2023/04/27
package lib

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/time"
	"note/src/util/validate"

	"strings"
)

// UpdName 图片重命名
func UpdName(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// lib
	img := typ.Lib{}
	err := api_common_context.ShouldBind(context, &img)
	if err != nil {
		redirect(err)
		return
	}

	// id
	id := img.Id
	if id <= 0 {
		redirect(err)
		return
	}

	// name
	name := strings.TrimSpace(img.Name)
	// validate name
	err = validate.FileName(name)
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = db.Upd(context, "UPDATE `lib` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.NowUnix(), id, name)

	// redirect
	redirect(err)
}
