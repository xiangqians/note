// note del
// @author xiangqian
// @date 22:57 2023/04/25
package note

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"note/app/api/common/db"
	typ2 "note/app/typ"
)

// Restore 恢复（还原）
func Restore(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// note
	note, count, err := DbQry(context, typ2.Note{Abs: typ2.Abs{Id: id, Del: 1}, Pid: -1})
	pid := note.Pid
	if err != nil || count == 0 {
		redirect(pid, err)
		return
	}

	// update
	_, err = db.DbUpd(context, "UPDATE `note` SET `del` = 0, `upd_time` = ? WHERE `id` = ?", time.NowUnix(), id)
	redirect(pid, err)
	return
}

// Del 删除文件
func Del(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// note
	note, count, err := DbQry(context, typ2.Note{Abs: typ2.Abs{Id: id}, Pid: -1})
	pid := note.Pid
	if err != nil || count == 0 {
		redirect(pid, err)
		return
	}

	// 如果是目录则校验目录下是否有子文件
	if typ2.ExtNameOf(note.Type) == typ2.FtD {
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
	_, err = db.DbDel(context, "UPDATE `note` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.NowUnix(), id)

	// redirect
	redirect(pid, err)
	return
}
