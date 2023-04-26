// note update
// @author xiangqian
// @date 22:58 2023/04/25
package note

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	util_time "note/src/util/time"
	"strings"
)

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// note
	note := typ_api.Note{}
	err := common.ShouldBind(context, &note)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = common.VerifyName(note.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", note.Name, util_time.NowUnix(), note.Id, note.Name)

	redirect(pid, err)
	return
}
