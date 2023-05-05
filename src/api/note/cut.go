// note cut
// @author xiangqian
// @date 22:57 2023/04/25
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// Cut 剪切文件
func Cut(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.ConvTypeToStr(err),
		}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", id), resp)
	}

	// dst id
	dstId, err := common.Param[int64](context, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := common.Param[int64](context, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		var note typ_api.Note
		var count int64
		note, count, err = DbQry(context, typ_api.Note{Abs: typ_api.Abs{Id: dstId}, Pid: -1})
		if err != nil || count == 0 || typ_ft.FtD != typ_ft.ExtNameOf(note.Type) { // 拷贝到目标类型必须是目录
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		util_time.NowUnix(),
		srcId,
		dstId)

	// redirect
	redirect(dstId, err)
	return
}
