// note update
// @author xiangqian
// @date 22:58 2023/04/25
package note

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	typ_ft "note/src/typ/ft"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
)

// UpdContent 修改文件内容
func UpdContent(context *gin.Context) {
	json := func(err error) {
		if err != nil {
			common.JsonBadRequest(context, typ.Resp[any]{Msg: util_str.ConvTypeToStr(err)})
			return
		}

		common.JsonOk(context, typ.Resp[any]{})
	}

	// id
	id, err := common.PostForm[int64](context, "id")
	if err != nil {
		json(err)
		return
	}
	//log.Println("id", id)

	// f
	f, count, err := DbQry(context, typ.Note{Abs: typ.Abs{Id: id}, Pid: -1})
	if count == 0 || typ_ft.ExtNameOf(f.Type) != typ_ft.FtMd {
		json(nil)
		return
	}

	// content
	content, err := common.PostForm[string](context, "content")
	if err != nil {
		json(err)
		return
	}
	//log.Println("content", content)

	// os file
	fPath, err := Path(context, f)
	if err != nil {
		json(err)
		return
	}
	pFile, err := os.OpenFile(fPath,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)
	if err != nil {
		json(err)
		return
	}
	defer pFile.Close()

	// write
	pWriter := bufio.NewWriter(pFile)
	pWriter.WriteString(content)
	pWriter.Flush()

	// file info
	fInfo, err := pFile.Stat()
	if err != nil {
		json(err)
		return
	}

	size := fInfo.Size()

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, util_time.NowUnix(), id)
	if err != nil {
		json(err)
		return
	}

	json(nil)
	return
}

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// note
	note := typ.Note{}
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
