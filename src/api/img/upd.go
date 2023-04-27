// img update
// @author xiangqian
// @date 21:58 2023/04/27
package img

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	util_time "note/src/util/time"
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

	// name
	img.Name = strings.TrimSpace(img.Name)
	err = common.VerifyName(img.Name)
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", img.Name, util_time.NowUnix(), img.Id, img.Name)

	// redirect
	redirect(err)
	return
}
