// img update
// @author xiangqian
// @date 21:58 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"strings"
)

// UpdName 图片重命名
func UpdName(context *gin.Context) {
	redirect := func(err any) {
		RedirectToList(context, err)
	}

	// img
	img := typ_api.Img{}
	err := common.ShouldBind(context, &img)
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
	err = util_validate.FileName(name)
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, util_time.NowUnix(), id, name)

	// redirect
	redirect(err)
}
