// note del
// @author xiangqian
// @date 22:57 2023/04/25
package note

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// Del 删除文件
func Del(context *gin.Context) {
	redirect := func(pid int64, err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// note
	note, _, err := DbQry(context, id, 0)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// 如果是目录则校验目录下是否有子文件
	if typ_ft.ExtNameOf(note.Type) == typ_ft.FtD {
		var count int64
		count, err = DbCount(context, id)
		if err != nil {
			redirect(pid, err)
			return
		}

		if count != 0 {
			redirect(pid, i18n.MustGetMessage("i18n.cannotDelNonEmptyDir"))
			return
		}
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `note` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", util_time.NowUnix(), id)

	// redirect
	redirect(pid, err)
	return
}
