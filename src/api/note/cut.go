// note cut
// @author xiangqian
// @date 22:57 2023/04/25
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/str"
	"note/src/util/time"
)

// Cut 剪切文件
func Cut(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ.Resp[any]{
			Msg: str.ConvTypeToStr(err),
		}
		api_common_context.Redirect(context, fmt.Sprintf("/note/list?pid=%d", id), resp)
	}

	// dst id
	dstId, err := api_common_context.Param[int64](context, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := api_common_context.Param[int64](context, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		var note typ.Note
		var count int64
		note, count, err = DbQry(context, typ.Note{Abs: typ.Abs{Id: dstId}, Pid: -1})
		if err != nil || count == 0 || typ.FtD != typ.ExtNameOf(note.Type) { // 拷贝到目标类型必须是目录
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		time.NowUnix(),
		srcId,
		dstId)

	// redirect
	redirect(dstId, err)
	return
}
