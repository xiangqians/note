// note update
// @author xiangqian
// @date 22:58 2023/04/25
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"strings"
)

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.TypeToStr(err),
		}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
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
