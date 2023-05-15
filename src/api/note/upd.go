// note update
// @author xiangqian
// @date 22:58 2023/04/25
package note

import (
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// note
	note := typ.Note{}
	err := api_common_context.ShouldBind(context, &note)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// id
	id := note.Id

	// name
	name := strings.TrimSpace(note.Name)
	// validate name
	err = validate.FileName(name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.NowUnix(), id, name)

	redirect(pid, err)
}
