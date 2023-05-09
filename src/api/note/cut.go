// note cut
// @author xiangqian
// @date 22:57 2023/04/25
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/app/api/common/db"
	typ2 "note/app/typ"
)

// Cut 剪切文件
func Cut(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ2.Resp[any]{
			Msg: str.ConvTypeToStr(err),
		}
		context.Redirect(context, fmt.Sprintf("/note/list?pid=%d", id), resp)
	}

	// dst id
	dstId, err := context.Param[int64](context, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := context.Param[int64](context, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		var note typ2.Note
		var count int64
		note, count, err = DbQry(context, typ2.Note{Abs: typ2.Abs{Id: dstId}, Pid: -1})
		if err != nil || count == 0 || typ2.FtD != typ2.ExtNameOf(note.Type) { // 拷贝到目标类型必须是目录
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = db.DbUpd(context, "UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		time.NowUnix(),
		srcId,
		dstId)

	// redirect
	redirect(dstId, err)
	return
}
